package enrichment

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/borkanie/parliament-scraper/internal/db"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/otiai10/gosseract/v2"
)

// LawEnricher fetches PDF documents for unenriched voting rounds and extracts
// title/description via text extraction or OCR, then enriches via OpenAI.
type LawEnricher struct {
	database *db.DB
	client   *http.Client
	openai   *openai.Client
	tmpDir   string
}

func NewLawEnricher(database *db.DB, openAIKey string) *LawEnricher {
	tmpDir, _ := os.MkdirTemp("", "parliament-scraper-*")
	var aiClient *openai.Client
	if openAIKey != "" {
		aiClient = openai.NewClient(option.WithAPIKey(openAIKey))
	}
	return &LawEnricher{
		database: database,
		client:   &http.Client{Timeout: 60 * time.Second},
		openai:   aiClient,
		tmpDir:   tmpDir,
		}
}

// EnrichPending fetches up to limit unenriched rounds and fills in title/description.
func (enricher *LawEnricher) EnrichPending(ctx context.Context, limit int) error {
	rounds, err := enricher.database.GetRoundsWithoutDescription(ctx, limit)
	if err != nil {
		return fmt.Errorf("get unenriched rounds: %w", err)
	}

	for _, round := range rounds {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := enricher.enrichRound(ctx, round.ID, round.VoteId); err != nil {
			slog.Warn("failed to enrich round", "voteId", round.VoteId, "err", err)
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (enricher *LawEnricher) enrichRound(ctx context.Context, roundID string, voteId int) error {
	// Try to find a project link from the nominal vote page
	htmlURL := fmt.Sprintf("https://www.cdep.ro/pls/steno/evot2015.Nominal?idv=%d", voteId)
	projectURL, err := enricher.findProjectURL(htmlURL)
	if err != nil || projectURL == "" {
		return nil // no document to enrich from
	}

	// Find PDF link from project page
	pdfURL, err := enricher.findPDFURL(projectURL)
	if err != nil || pdfURL == "" {
		return nil
	}

	text, err := enricher.extractTextFromPDF(ctx, pdfURL)
	if err != nil || strings.TrimSpace(text) == "" {
		return nil
	}

	if enricher.openai == nil {
		slog.Warn("no OpenAI key set, skipping description enrichment", "voteId", voteId)
		return nil
	}

	title, description, err := enricher.callOpenAI(ctx, text)
	if err != nil {
		return fmt.Errorf("openai call: %w", err)
	}

	return enricher.database.UpdateRoundDescription(ctx, roundID, title, description)
}

func (enricher *LawEnricher) findProjectURL(htmlURL string) (string, error) {
	resp, err := enricher.client.Get(htmlURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`href="(/pls/proiecte/upl_pck2015\.proiect\?[^"]+)"`)
	match := re.FindSubmatch(body)
	if match == nil {
		return "", nil
	}
	return "https://www.cdep.ro" + string(match[1]), nil
}

func (enricher *LawEnricher) findPDFURL(projectURL string) (string, error) {
	resp, err := enricher.client.Get(projectURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Try promulgation, then adopted form, then motivation in order
	patterns := []string{
		`href="((?:/docs|/proiecte)[^"]*_pr\.pdf)"`,
		`href="((?:/proiecte)[^"]*(?:se|cd)\d+\.pdf)"`,
		`href="((?:/proiecte)[^"]*em\d+\.pdf)"`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindSubmatch(body); match != nil {
			parsedURL := string(match[1])
			if !strings.HasPrefix(parsedURL, "http") {
				parsedURL = "https://www.cdep.ro" + parsedURL
			}
			return parsedURL, nil
		}
	}
	return "", nil
}

func (enricher *LawEnricher) extractTextFromPDF(ctx context.Context, pdfURL string) (string, error) {
	// Download PDF to temp file
	resp, err := enricher.client.Get(pdfURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile := filepath.Join(enricher.tmpDir, fmt.Sprintf("law_%d.pdf", time.Now().UnixNano()))
	fileHandle, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile)

	if _, err := io.Copy(fileHandle, resp.Body); err != nil {
		fileHandle.Close()
		return "", err
	}
	fileHandle.Close()

	// Try pdftotext first (for selectable text PDFs)
	if text, err := pdfToText(tmpFile); err == nil && strings.TrimSpace(text) != "" {
		return text, nil
	}

	// Fall back to OCR
	return enricher.ocrPDF(tmpFile)
}

// pdfToText uses the poppler pdftotext utility for selectable text.
func pdfToText(pdfPath string) (string, error) {
	out, err := exec.Command("pdftotext", "-enc", "UTF-8", pdfPath, "-").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ocrPDF converts each PDF page to an image and runs Tesseract OCR.
func (enricher *LawEnricher) ocrPDF(pdfPath string) (string, error) {
	// Convert PDF pages to PNG images using pdftoppm
	imgPrefix := filepath.Join(enricher.tmpDir, fmt.Sprintf("page_%d", time.Now().UnixNano()))
	if err := exec.Command("pdftoppm", "-r", "300", "-png", pdfPath, imgPrefix).Run(); err != nil {
		return "", fmt.Errorf("pdftoppm: %w", err)
	}

	// Read all generated page files
	pageImages, err := filepath.Glob(imgPrefix + "-*.png")
	if err != nil || len(pageImages) == 0 {
		return "", fmt.Errorf("no page images generated")
	}

	ocrClient := gosseract.NewClient()
	defer ocrClient.Close()
	_ = ocrClient.SetLanguage("ron")

	var builder strings.Builder
	for _, img := range pageImages {
		_ = ocrClient.SetImage(img)
		text, err := ocrClient.Text()
		if err == nil {
			builder.WriteString(text)
			builder.WriteString("\n")
		}
		os.Remove(img)
	}

	return builder.String(), nil
}

func (enricher *LawEnricher) callOpenAI(ctx context.Context, lawText string) (title, description string, err error) {
	truncated := lawText
	if len(truncated) > 4000 {
		truncated = truncated[:4000]
	}

	prompt := fmt.Sprintf(`Ești un asistent juridic român. Din textul legii de mai jos extrage:
1. titlu - titlul scurt al legii (max 100 caractere)
2. descriere - o descriere succintă a legii (max 300 caractere)

Răspunde STRICT în formatul:
titlu: <titlul legii>
descriere: <descrierea legii>

Text lege:
%s`, truncated)

	resp, err := enricher.openai.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
	})
	if err != nil {
		return "", "", err
	}

	if len(resp.Choices) == 0 {
		return "", "", fmt.Errorf("empty response from OpenAI")
	}

	return parseOpenAIResponse(resp.Choices[0].Message.Content)
}

func parseOpenAIResponse(content string) (title, description string, err error) {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "titlu:"); ok {
			title = strings.TrimSpace(after)
		}
		if after, ok := strings.CutPrefix(line, "descriere:"); ok {
			description = strings.TrimSpace(after)
		}
	}
	if title == "" {
		return "", "", fmt.Errorf("could not parse title from OpenAI response")
	}
	return title, description, nil
}

// Cleanup removes the temporary directory used for PDF processing.
func (enricher *LawEnricher) Cleanup() {
	os.RemoveAll(enricher.tmpDir)
}
