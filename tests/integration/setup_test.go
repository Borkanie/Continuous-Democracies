package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/borkanie/parliament-api/internal/db"
	"github.com/borkanie/parliament-api/internal/handlers"
	"github.com/borkanie/parliament-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tcmongo "github.com/testcontainers/testcontainers-go/modules/mongodb"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start MongoDB container
	mongoC, err := tcmongo.Run(ctx, "mongo:7")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start MongoDB container: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = mongoC.Terminate(ctx) }()

	uri, err := mongoC.ConnectionString(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get MongoDB URI: %v\n", err)
		os.Exit(1)
	}

	database, err := db.Connect(uri, "testdb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to test MongoDB: %v\n", err)
		os.Exit(1)
	}

	idxCtx, idxCancel := context.WithTimeout(ctx, 30*time.Second)
	defer idxCancel()
	if err := db.EnsureIndexes(idxCtx, database); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create indexes: %v\n", err)
		os.Exit(1)
	}

	if err := seedDatabase(ctx, database); err != nil {
		fmt.Fprintf(os.Stderr, "failed to seed database: %v\n", err)
		os.Exit(1)
	}

	// Build the router (same wiring as main.go)
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	partyRepo := repository.NewPartyRepo(database)
	politiciansRepo := repository.NewPoliticiansRepo(database)
	votingRepo := repository.NewVotingRepo(database)

	partyH := handlers.NewPartyHandler(partyRepo)
	politiciansH := handlers.NewPoliticiansHandler(politiciansRepo)
	votingH := handlers.NewVotingHandler(votingRepo)

	router.Get("/health", handlers.Health)
	router.Route("/api/Party", func(router chi.Router) {
		router.Get("/all", partyH.All)
		router.Get("/GetById/", partyH.GetByID)
		router.Get("/query", partyH.Query)
	})
	router.Route("/api/Politicians", func(router chi.Router) {
		router.Get("/getAllPoliticians", politiciansH.GetAll)
		router.Get("/GetById/", politiciansH.GetByID)
		router.Get("/GetByName/", politiciansH.GetByName)
	})
	router.Route("/api/Voting", func(router chi.Router) {
		router.Get("/getAllRounds", votingH.GetAllRounds)
		router.Get("/getRoundById/", votingH.GetRoundById)
		router.Get("/GetResultForVote/", votingH.GetResultForVote)
		router.Get("/GetAllVotesForARoundById/", votingH.GetAllVotesForARoundById)
	})

	testServer = httptest.NewServer(router)
	defer testServer.Close()

	os.Exit(m.Run())
}

func seedDatabase(ctx context.Context, database *mongodriver.Database) error {
	fixtures := []struct {
		collection string
		file       string
	}{
		{db.ColParties, "../../tests/fixtures/parties.json"},
		{db.ColPoliticians, "../../tests/fixtures/politicians.json"},
		{db.ColVotingRounds, "../../tests/fixtures/voting_rounds.json"},
		{db.ColVotes, "../../tests/fixtures/votes.json"},
	}

	for _, fixture := range fixtures {
		data, err := os.ReadFile(fixture.file)
		if err != nil {
			return fmt.Errorf("read fixture %s: %w", fixture.file, err)
		}

		var documents []bson.M
		if err := json.Unmarshal(data, &documents); err != nil {
			return fmt.Errorf("parse fixture %s: %w", fixture.file, err)
		}

		col := database.Collection(fixture.collection)
		docsInterface := make([]interface{}, len(documents))
		for index, document := range documents {
			docsInterface[index] = document
		}

		_, err = col.InsertMany(ctx, docsInterface, options.InsertMany())
		if err != nil {
			return fmt.Errorf("insert fixture %s: %w", fixture.file, err)
		}
	}
	return nil
}
