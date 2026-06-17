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

const (
	ColParties      = "parties"
	ColPoliticians  = "politicians"
	ColVotingRounds = "voting_rounds"
	ColVotes        = "votes"
)

func Connect(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	slog.Info("connected to MongoDB", "uri", uri, "db", dbName)
	return client.Database(dbName), nil
}

func EnsureIndexes(ctx context.Context, database *mongo.Database) error {
	type indexSpec struct {
		collection string
		models     []mongo.IndexModel
	}

	specs := []indexSpec{
		{
			collection: ColParties,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "acronym", Value: 1}}, Options: options.Index().SetSparse(true)},
				{Keys: bson.D{{Key: "name", Value: 1}}},
			},
		},
		{
			collection: ColPoliticians,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "partyId", Value: 1}}},
				{Keys: bson.D{{Key: "name", Value: 1}}},
				{Keys: bson.D{{Key: "active", Value: 1}}},
			},
		},
		{
			collection: ColVotingRounds,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "voteId", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "voteDate", Value: -1}}},
				{
					Keys: bson.D{
						{Key: "title", Value: "text"},
						{Key: "description", Value: "text"},
					},
				},
			},
		},
		{
			collection: ColVotes,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "roundId", Value: 1}}},
				{Keys: bson.D{{Key: "roundId", Value: 1}, {Key: "position", Value: 1}}},
				{Keys: bson.D{{Key: "politicianId", Value: 1}}},
				{
					Keys:    bson.D{{Key: "politicianId", Value: 1}, {Key: "roundId", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	}

	for _, spec := range specs {
		collection := database.Collection(spec.collection)
		if _, err := collection.Indexes().CreateMany(ctx, spec.models); err != nil {
			return fmt.Errorf("create indexes for %s: %w", spec.collection, err)
		}
	}

	slog.Info("database indexes ensured")
	return nil
}
