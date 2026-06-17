package importer

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/borkanie/parliament-scraper/internal/db"
	"github.com/borkanie/parliament-scraper/internal/models"
	"github.com/google/uuid"
)

const (
	cdepBase       = "https://www.cdep.ro"
	defaultStartId = 35333
)

// xmlVoteRow mirrors the XML structure from evot2015.xml
type xmlVoteRow struct {
	VoteID  string `xml:"VOTID"`
	Surname string `xml:"NUME"`
	Name    string `xml:"PRENUME"`
	Group   string `xml:"GRUP"`
	Camera  string `xml:"CAMERA"`
	Vote    string `xml:"VOT"`
}

type xmlRoot struct {
	Rows []xmlVoteRow `xml:"ROW"`
}

// VotesImporter handles vote data ingestion from cdep.ro.
type VotesImporter struct {
	db     *db.DB
	client *http.Client
}

func NewVotesImporter(database *db.DB, proxyURL string) *VotesImporter {
	transport := http.DefaultTransport
	if proxyURL != "" {
		transport = socksTransport(proxyURL)
	}
	return &VotesImporter{
		db:     database,
		client: &http.Client{Transport: transport, Timeout: 30 * time.Second},
	}
}

// GoFromLastForward imports up to maxForward new voting rounds starting from the DB max.
func (importer *VotesImporter) GoFromLastForward(ctx context.Context, maxForward int) error {
	dbMax, err := importer.db.GetMaxVoteId(ctx)
	if err != nil {
		return fmt.Errorf("get max voteId: %w", err)
	}
	if dbMax == 0 {
		dbMax = defaultStartId
	}

	remoteMax, err := importer.fetchRemoteMaxVoteId()
	if err != nil {
		slog.Warn("could not fetch remote max voteId, falling back to forward scan", "err", err)
	}

	start := dbMax + 1
	end := dbMax + maxForward
	if remoteMax > dbMax {
		end = remoteMax
		slog.Info("remote has newer votes, syncing", "from", start, "to", end)
	} else {
		slog.Info("scanning forward", "from", start, "to", end)
	}

	return importer.importRange(ctx, start, end)
}

func (importer *VotesImporter) importRange(ctx context.Context, start, end int) error {
	for lawId := start; lawId <= end; lawId++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := importer.importLaw(ctx, lawId); err != nil {
			slog.Error("failed to import law", "lawId", lawId, "err", err)
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func (importer *VotesImporter) importLaw(ctx context.Context, lawId int) error {
	xmlURL := fmt.Sprintf("%s/pls/steno/evot2015.xml?par1=2&par2=%d", cdepBase, lawId)
	rows, err := importer.fetchVoteXML(xmlURL)
	if err != nil {
		return fmt.Errorf("fetch xml: %w", err)
	}
	if len(rows) == 0 {
		return nil // no votes for this law
	}

	// Get or create the voting round
	round, err := importer.db.GetRoundByVoteId(ctx, lawId)
	if err != nil {
		return fmt.Errorf("get round: %w", err)
	}
	if round == nil {
		round, err = importer.extractAndInsertRound(ctx, lawId)
		if err != nil {
			return fmt.Errorf("extract round: %w", err)
		}
	}

	imported := 0
	for _, row := range rows {
		if err := importer.processVoteRow(ctx, row, round.ID); err != nil {
			slog.Warn("skipping vote row", "lawId", lawId, "name", row.Surname, "err", err)
			continue
		}
		imported++
	}

	slog.Info("imported votes", "lawId", lawId, "total", imported)
	return nil
}

func (importer *VotesImporter) processVoteRow(ctx context.Context, row xmlVoteRow, roundId string) error {
	group := row.Group
	if group == "Minoritati" {
		group = "MIN"
	}

	// Ensure party exists
	party, err := importer.db.GetPartyByAcronym(ctx, group)
	if err != nil {
		return err
	}
	if party == nil {
		newParty := &models.Party{
			ID:      uuid.NewString(),
			Acronym: group,
			Name:    group,
			Color:   "#727272",
			Active:  true,
		}
		if err := importer.db.UpsertParty(ctx, newParty); err != nil {
			return fmt.Errorf("upsert party %s: %w", group, err)
		}
		party, err = importer.db.GetPartyByAcronym(ctx, group)
		if err != nil || party == nil {
			return fmt.Errorf("party not found after insert: %s", group)
		}
	}

	fullName := normalizeFullName(row.Name, row.Surname)

	// Ensure politician exists
	politician, err := importer.db.GetPoliticianByNameAndParty(ctx, fullName, party.ID)
	if err != nil {
		return err
	}
	if politician == nil {
		newPolitician := &models.Politician{
			ID:      uuid.NewString(),
			Name:    fullName,
			Gender:  inferGender(row.Name),
			PartyId: party.ID,
			Active:  true,
		}
		if err := importer.db.UpsertPolitician(ctx, newPolitician); err != nil {
			return fmt.Errorf("upsert politician %s: %w", fullName, err)
		}
		politician, err = importer.db.GetPoliticianByNameAndParty(ctx, fullName, party.ID)
		if err != nil || politician == nil {
			return fmt.Errorf("politician not found after insert: %s", fullName)
		}
	}

	vote := &models.Vote{
		ID:           uuid.NewString(),
		PoliticianId: politician.ID,
		RoundId:      roundId,
		Position:     votePosition(row.Vote),
	}
	return importer.db.InsertVoteIfNotExists(ctx, vote)
}

func (importer *VotesImporter) extractAndInsertRound(ctx context.Context, lawId int) (*models.VotingRound, error) {
	htmlURL := fmt.Sprintf("%s/pls/steno/evot2015.Nominal?idv=%d", cdepBase, lawId)
	voteDate, err := importer.fetchVoteDate(htmlURL)
	if err != nil {
		slog.Warn("could not extract vote date, using now", "lawId", lawId, "err", err)
		voteDate = time.Now().UTC()
	}

	title := fmt.Sprintf("Vot electronic:%d", lawId)
	round := &models.VotingRound{
		ID:       uuid.NewString(),
		VoteId:   lawId,
		Title:    title,
		Name:     title,
		VoteDate: voteDate,
	}
	if err := importer.db.InsertRoundIfNotExists(ctx, round); err != nil {
		return nil, err
	}
	return importer.db.GetRoundByVoteId(ctx, lawId)
}

func (importer *VotesImporter) fetchVoteXML(xmlURL string) ([]xmlVoteRow, error) {
	resp, err := importer.client.Get(xmlURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var root xmlRoot
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	return root.Rows, nil
}

func (importer *VotesImporter) fetchVoteDate(htmlURL string) (time.Time, error) {
	resp, err := importer.client.Get(htmlURL)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return time.Time{}, err
	}

	// The date appears in a td with class "vot-title" or as text in the page
	text := doc.Find("td.vot-title").First().Text()
	if text == "" {
		text = doc.Find("td").FilterFunction(func(_ int, selection *goquery.Selection) bool {
			return strings.Contains(selection.Text(), "Data votului")
		}).Next().Text()
	}

	return parseCdepDate(strings.TrimSpace(text))
}

func (importer *VotesImporter) fetchRemoteMaxVoteId() (int, error) {
	resp, err := importer.client.Get(cdepBase + "/pls/steno/evot2015.data")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(`[?&]idv=(\d+)`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	maxVoteId := 0
	for _, match := range matches {
		if number, err := strconv.Atoi(match[1]); err == nil && number > maxVoteId {
			maxVoteId = number
		}
	}
	if maxVoteId == 0 {
		return 0, fmt.Errorf("no idv= entries found on data page")
	}
	return maxVoteId, nil
}

// normalizeFullName converts "PRENUME NUME" (first last) to "Prenume Nume"
func normalizeFullName(firstName, surname string) string {
	parts := strings.Fields(strings.ToLower(firstName + " " + surname))
	for index, part := range parts {
		if len(part) > 0 {
			parts[index] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

// inferGender uses a heuristic: Romanian female first names typically end in 'a'
func inferGender(firstName string) int {
	if strings.HasSuffix(strings.ToLower(strings.TrimSpace(firstName)), "a") {
		return 1 // female
	}
	return 0 // male
}

// votePosition maps Romanian vote strings to position integers
func votePosition(vot string) int {
	switch strings.TrimSpace(strings.ToUpper(vot)) {
	case "DA":
		return models.PositionYes
	case "NU":
		return models.PositionNo
	case "AB":
		return models.PositionAbstain
	default:
		return models.PositionAbsent
	}
}

// parseCdepDate tries to parse date strings found on cdep.ro pages
func parseCdepDate(dateString string) (time.Time, error) {
	formats := []string{
		"02.01.2006 15:04:05",
		"02.01.2006 15:04",
		"02.01.2006",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
	}
	for _, format := range formats {
		if parsedTime, err := time.Parse(format, dateString); err == nil {
			return parsedTime.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse date: %q", dateString)
}
