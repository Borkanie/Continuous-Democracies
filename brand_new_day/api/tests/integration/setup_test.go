package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/borkanie/brand-new-day-api/internal/controller"
	"github.com/borkanie/brand-new-day-api/internal/db"
	"github.com/borkanie/brand-new-day-api/internal/generated"
	"github.com/borkanie/brand-new-day-api/internal/models"
	"github.com/borkanie/brand-new-day-api/internal/repository"
	"github.com/borkanie/brand-new-day-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tcmongo "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var testServer *httptest.Server

func TestMain(testingManager *testing.M) {
	testContext := context.Background()

	// Start MongoDB container
	mongoContainer, err := tcmongo.Run(testContext, "mongo:7")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start MongoDB container: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = mongoContainer.Terminate(testContext) }()

	mongoConnectionString, err := mongoContainer.ConnectionString(testContext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get MongoDB URI: %v\n", err)
		os.Exit(1)
	}

	testDatabase, err := db.Connect(mongoConnectionString, "testdb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to test MongoDB: %v\n", err)
		os.Exit(1)
	}

	indexContext, cancelIndexContext := context.WithTimeout(testContext, 30*time.Second)
	defer cancelIndexContext()
	if err := db.EnsureIndexes(indexContext, testDatabase); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create indexes: %v\n", err)
		os.Exit(1)
	}

	if err := seedTestDatabase(testContext, testDatabase); err != nil {
		fmt.Fprintf(os.Stderr, "failed to seed database: %v\n", err)
		os.Exit(1)
	}

	// Build the router (same wiring as main.go)
	testRouter := chi.NewRouter()
	testRouter.Use(middleware.Recoverer)

	partyRepository := repository.NewPartyRepository(testDatabase)
	politicianRepository := repository.NewPoliticianRepository(testDatabase)
	lawBucketRepository := repository.NewLawBucketRepository(testDatabase)
	normativeRepository := repository.NewNormativeRepository(testDatabase)
	votingRoundRepository := repository.NewVotingRoundRepository(testDatabase)

	politicianService := service.NewPoliticianService(politicianRepository, partyRepository, votingRoundRepository)
	votingService := service.NewVotingService(
		votingRoundRepository,
		normativeRepository,
		lawBucketRepository,
		politicianRepository,
		partyRepository,
	)
	lawService := service.NewLawService(lawBucketRepository, normativeRepository, votingRoundRepository)

	apiController := controller.NewController(politicianService, votingService, lawService)

	generated.HandlerFromMux(apiController, testRouter)

	testServer = httptest.NewServer(testRouter)
	defer testServer.Close()

	exitCode := testingManager.Run()
	os.Exit(exitCode)
}

func seedTestDatabase(seedContext context.Context, testDatabase *mongodriver.Database) error {
	fixtureList := []struct {
		collectionName  string
		fixtureFilePath string
		unmarshalTarget any
	}{
		{
			collectionName:  db.CollectionParties,
			fixtureFilePath: "../fixtures/parties.json",
			unmarshalTarget: &[]bson.M{},
		},
		{
			collectionName:  db.CollectionPoliticians,
			fixtureFilePath: "../fixtures/politicians.json",
			unmarshalTarget: &[]bson.M{},
		},
		{
			collectionName:  db.CollectionLawBuckets,
			fixtureFilePath: "../fixtures/law_buckets.json",
			unmarshalTarget: &[]bson.M{},
		},
		{
			collectionName:  db.CollectionNormatives,
			fixtureFilePath: "../fixtures/normatives.json",
			unmarshalTarget: &[]interface{}{},
		},
		{
			collectionName:  db.CollectionVotingRounds,
			fixtureFilePath: "../fixtures/voting_rounds.json",
			unmarshalTarget: &[]bson.M{},
		},
	}

	for _, fixtureEntry := range fixtureList {
		fixtureData, err := os.ReadFile(fixtureEntry.fixtureFilePath)
		if err != nil {
			return fmt.Errorf("read fixture %s: %w", fixtureEntry.fixtureFilePath, err)
		}

		// For normatives, we need special handling because of the compound _id
		if fixtureEntry.collectionName == db.CollectionNormatives {
			var normativeDocuments []struct {
				ID            models.NormativeKey `json:"_id" bson:"_id"`
				LawBucketID   int                 `json:"lawBucketId" bson:"lawBucketId"`
				Type          string              `json:"type" bson:"type"`
				Label         string              `json:"label" bson:"label"`
				Text          string              `json:"text" bson:"text"`
				EffectiveDate time.Time           `json:"effectiveDate" bson:"effectiveDate"`
			}

			if err := json.Unmarshal(fixtureData, &normativeDocuments); err != nil {
				return fmt.Errorf("parse normatives fixture: %w", err)
			}

			documentInterfaces := make([]interface{}, len(normativeDocuments))
			for documentIndex, normativeDocument := range normativeDocuments {
				documentInterfaces[documentIndex] = normativeDocument
			}

			collection := testDatabase.Collection(fixtureEntry.collectionName)
			_, err = collection.InsertMany(seedContext, documentInterfaces, options.InsertMany())
			if err != nil {
				return fmt.Errorf("insert fixture %s: %w", fixtureEntry.fixtureFilePath, err)
			}
		} else {
			var documents []bson.M
			if err := json.Unmarshal(fixtureData, &documents); err != nil {
				return fmt.Errorf("parse fixture %s: %w", fixtureEntry.fixtureFilePath, err)
			}

			documentInterfaces := make([]interface{}, len(documents))
			for documentIndex, document := range documents {
				documentInterfaces[documentIndex] = document
			}

			collection := testDatabase.Collection(fixtureEntry.collectionName)
			_, err = collection.InsertMany(seedContext, documentInterfaces, options.InsertMany())
			if err != nil {
				return fmt.Errorf("insert fixture %s: %w", fixtureEntry.fixtureFilePath, err)
			}
		}
	}

	return nil
}
