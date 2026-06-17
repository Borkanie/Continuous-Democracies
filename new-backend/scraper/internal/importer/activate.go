package importer

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/borkanie/parliament-scraper/internal/db"
)

// ActivationImporter refreshes which politicians are currently active in Parliament.
type ActivationImporter struct {
	db     *db.DB
	client *http.Client
}

func NewActivationImporter(database *db.DB, proxyURL string) *ActivationImporter {
	transport := http.DefaultTransport
	if proxyURL != "" {
		transport = socksTransport(proxyURL)
	}
	return &ActivationImporter{
		db:     database,
		client: &http.Client{Transport: transport, Timeout: 30 * time.Second},
	}
}

// ActivateAllGroups fetches current parliamentary groups from cdep.ro and
// updates Active status for all politicians in the database.
func (importer *ActivationImporter) ActivateAllGroups(ctx context.Context) error {
	const idl = 1
	var allNames []string

	for idg := 1; ; idg++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		groupURL := fmt.Sprintf("https://www.cdep.ro/pls/parlam/structura2015.gp?idl=%d&idg=%d", idl, idg)
		names, found, err := importer.fetchGroupNames(groupURL)
		if err != nil {
			slog.Warn("failed to fetch group page", "idg", idg, "err", err)
			break
		}
		if !found {
			break
		}
		allNames = append(allNames, names...)
		slog.Debug("fetched group", "idg", idg, "names", len(names))
	}

	// Deduplicate
	seen := make(map[string]bool, len(allNames))
	unique := make([]string, 0, len(allNames))
	for _, name := range allNames {
		if !seen[name] {
			seen[name] = true
			unique = append(unique, name)
		}
	}

	slog.Info("activating politicians", "count", len(unique))
	return importer.db.ReplaceActivePoliticians(ctx, unique)
}

func (importer *ActivationImporter) fetchGroupNames(groupURL string) (names []string, found bool, err error) {
	resp, err := importer.client.Get(groupURL)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, false, err
	}

	grp := doc.Find("div.grup-parlamentar-list.grupuri-parlamentare-list")
	if grp.Length() == 0 {
		return nil, false, nil // stop iteration
	}

	grp.Find("a").Each(func(_ int, selection *goquery.Selection) {
		text := strings.TrimSpace(selection.Text())
		if text != "" {
			names = append(names, moveFirstWordToEnd(text))
		}
	})

	return names, true, nil
}

// moveFirstWordToEnd converts "POPESCU Ion" (Surname First) → "Ion POPESCU"
// cdep.ro displays names as "SURNAME Firstname"; the DB stores "Firstname SURNAME"
func moveFirstWordToEnd(nameText string) string {
	parts := strings.Fields(nameText)
	if len(parts) <= 1 {
		return nameText
	}
	return strings.Join(append(parts[1:], parts[0]), " ")
}
