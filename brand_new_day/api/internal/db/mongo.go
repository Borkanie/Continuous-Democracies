// Package db owns the MongoDB connection and index setup for
// brand_new_day/api. Pattern mirrors new-backend/api/internal/db/mongo.go,
// adapted to this module's collections and to the v2 mongo driver.
package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Collection name constants. These are the cross-worker contract shared with
// the seeder and every repository in internal/repository.
const (
	CollectionParties      = "parties"
	CollectionPoliticians  = "politicians"
	CollectionLawBuckets   = "law_buckets"
	CollectionNormatives   = "normatives"
	CollectionVotingRounds = "voting_rounds"
)

// connectTimeout bounds how long the initial connect+ping is allowed to take.
const connectTimeout = 10 * time.Second

// ConnectClient dials MongoDB at mongoURI, verifies the connection with a
// ping, and returns both the live *mongo.Client (so the caller can
// Disconnect it cleanly on shutdown) and the *mongo.Database handle for
// databaseName.
func ConnectClient(mongoURI string, databaseName string) (*mongo.Client, *mongo.Database, error) {
	connectContext, cancelConnect := context.WithTimeout(context.Background(), connectTimeout)
	defer cancelConnect()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, nil, fmt.Errorf("mongo connect: %w", err)
	}

	if err := client.Ping(connectContext, nil); err != nil {
		return nil, nil, fmt.Errorf("mongo ping: %w", err)
	}

	slog.Info("connected to MongoDB", "uri", mongoURI, "database", databaseName)
	return client, client.Database(databaseName), nil
}

// Connect is a convenience wrapper around ConnectClient for callers that
// don't need the underlying *mongo.Client (e.g. short-lived scripts). Callers
// that need a clean shutdown path (the long-running API server) should use
// ConnectClient instead so they can Disconnect the client on exit.
func Connect(mongoURI string, databaseName string) (*mongo.Database, error) {
	_, database, err := ConnectClient(mongoURI, databaseName)
	if err != nil {
		return nil, err
	}
	return database, nil
}

// EnsureIndexes creates every index required by the data model described in
// ARCHITECTURE.md. It is idempotent: CreateOne/CreateMany against an
// already-existing equivalent index is a no-op.
func EnsureIndexes(requestContext context.Context, database *mongo.Database) error {
	type indexSpec struct {
		collectionName string
		indexModels    []mongo.IndexModel
	}

	specs := []indexSpec{
		{
			collectionName: CollectionParties,
			indexModels: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "acronym", Value: 1}},
					Options: options.Index().SetUnique(true).SetSparse(true),
				},
			},
		},
		{
			collectionName: CollectionNormatives,
			indexModels: []mongo.IndexModel{
				{Keys: bson.D{{Key: "lawBucketId", Value: 1}}},
			},
		},
		{
			collectionName: CollectionVotingRounds,
			indexModels: []mongo.IndexModel{
				// Multikey: indexing an array field (votes.politicianId)
				// automatically produces a multikey index. Serves "all votes
				// by this politician".
				{Keys: bson.D{{Key: "votes.politicianId", Value: 1}}},
				{Keys: bson.D{{Key: "normativeId", Value: 1}, {Key: "normativeVersion", Value: 1}}},
				{Keys: bson.D{{Key: "chamber", Value: 1}}},
				{
					Keys: bson.D{
						{Key: "title", Value: "text"},
						{Key: "description", Value: "text"},
					},
				},
			},
		},
	}

	for _, spec := range specs {
		collection := database.Collection(spec.collectionName)
		if _, err := collection.Indexes().CreateMany(requestContext, spec.indexModels); err != nil {
			return fmt.Errorf("create indexes for %s: %w", spec.collectionName, err)
		}
	}

	slog.Info("database indexes ensured")
	return nil
}
