package orchestrator

import (
	"context"
	"log/slog"
	"time"

	"github.com/borkanie/parliament-scraper/internal/db"
	"github.com/borkanie/parliament-scraper/internal/enrichment"
	"github.com/borkanie/parliament-scraper/internal/importer"
)

const (
	cycleInterval  = 12 * time.Hour
	forwardWindow  = 100
	enrichBatch    = 50
)

// Orchestrator coordinates the vote importer, law enricher, and politician activator.
type Orchestrator struct {
	votes      *importer.VotesImporter
	activation *importer.ActivationImporter
	enricher   *enrichment.LawEnricher
}

func New(database *db.DB, openAIKey, proxyURL string) *Orchestrator {
	return &Orchestrator{
		votes:      importer.NewVotesImporter(database, proxyURL),
		activation: importer.NewActivationImporter(database, proxyURL),
		enricher:   enrichment.NewLawEnricher(database, openAIKey),
	}
}

// RunOnce executes one full cycle: import new votes, enrich laws, refresh active status.
func (orchestrator *Orchestrator) RunOnce(ctx context.Context) {
	slog.Info("cycle start")

	slog.Info("importing new votes")
	if err := orchestrator.votes.GoFromLastForward(ctx, forwardWindow); err != nil {
		slog.Error("vote import failed", "err", err)
	}

	slog.Info("enriching law descriptions")
	if err := orchestrator.enricher.EnrichPending(ctx, enrichBatch); err != nil {
		slog.Error("law enrichment failed", "err", err)
	}

	slog.Info("refreshing active politicians")
	if err := orchestrator.activation.ActivateAllGroups(ctx); err != nil {
		slog.Error("politician activation failed", "err", err)
	}

	slog.Info("cycle complete")
}

// Run executes cycles on a 12-hour ticker until ctx is cancelled.
// Call RunOnce instead when running as a Kubernetes CronJob (RUN_ONCE=true).
func (orchestrator *Orchestrator) Run(ctx context.Context) {
	orchestrator.RunOnce(ctx)

	ticker := time.NewTicker(cycleInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("orchestrator shutting down")
			orchestrator.enricher.Cleanup()
			return
		case <-ticker.C:
			orchestrator.RunOnce(ctx)
		}
	}
}
