package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/borkanie/parliament-scraper/internal/models"
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

type DB struct {
	database *mongo.Database
}

func Connect(uri, dbName string) (*DB, error) {
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
	return &DB{database: client.Database(dbName)}, nil
}

// UpsertParty inserts or updates a party by acronym.
func (database *DB) UpsertParty(ctx context.Context, party *models.Party) error {
	col := database.database.Collection(ColParties)
	filter := bson.D{{Key: "acronym", Value: party.Acronym}}
	update := bson.D{{Key: "$setOnInsert", Value: bson.D{{Key: "_id", Value: party.ID}}}, {Key: "$set", Value: bson.D{
		{Key: "name", Value: party.Name},
		{Key: "acronym", Value: party.Acronym},
		{Key: "color", Value: party.Color},
		{Key: "active", Value: party.Active},
	}}}
	_, err := col.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))
	return err
}

// GetPartyByAcronym returns a party by acronym, or nil if not found.
func (database *DB) GetPartyByAcronym(ctx context.Context, acronym string) (*models.Party, error) {
	col := database.database.Collection(ColParties)
	var party models.Party
	err := col.FindOne(ctx, bson.D{{Key: "acronym", Value: acronym}}).Decode(&party)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &party, err
}

// UpsertPolitician inserts or updates a politician by name + partyId.
func (database *DB) UpsertPolitician(ctx context.Context, politician *models.Politician) error {
	col := database.database.Collection(ColPoliticians)
	filter := bson.D{{Key: "name", Value: politician.Name}, {Key: "partyId", Value: politician.PartyId}}
	update := bson.D{{Key: "$setOnInsert", Value: bson.D{{Key: "_id", Value: politician.ID}}}, {Key: "$set", Value: bson.D{
		{Key: "name", Value: politician.Name},
		{Key: "gender", Value: politician.Gender},
		{Key: "imageUrl", Value: politician.ImageUrl},
		{Key: "partyId", Value: politician.PartyId},
		{Key: "active", Value: politician.Active},
		{Key: "workLocation", Value: politician.WorkLocation},
	}}}
	_, err := col.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))
	return err
}

// GetPoliticianByNameAndParty returns a politician by name + partyId, or nil if not found.
func (database *DB) GetPoliticianByNameAndParty(ctx context.Context, name, partyId string) (*models.Politician, error) {
	col := database.database.Collection(ColPoliticians)
	var politician models.Politician
	err := col.FindOne(ctx, bson.D{{Key: "name", Value: name}, {Key: "partyId", Value: partyId}}).Decode(&politician)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &politician, err
}

// InsertRoundIfNotExists inserts a voting round; no-ops if voteId already exists.
func (database *DB) InsertRoundIfNotExists(ctx context.Context, round *models.VotingRound) error {
	col := database.database.Collection(ColVotingRounds)
	_, err := col.UpdateOne(
		ctx,
		bson.D{{Key: "voteId", Value: round.VoteId}},
		bson.D{{Key: "$setOnInsert", Value: round}},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}

// GetRoundByVoteId returns a round by its integer voteId, or nil if not found.
func (database *DB) GetRoundByVoteId(ctx context.Context, voteId int) (*models.VotingRound, error) {
	col := database.database.Collection(ColVotingRounds)
	var round models.VotingRound
	err := col.FindOne(ctx, bson.D{{Key: "voteId", Value: voteId}}).Decode(&round)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &round, err
}

// GetMaxVoteId returns the highest voteId currently in the DB, or 0 if empty.
func (database *DB) GetMaxVoteId(ctx context.Context) (int, error) {
	col := database.database.Collection(ColVotingRounds)
	opts := options.FindOne().SetSort(bson.D{{Key: "voteId", Value: -1}})
	var round models.VotingRound
	err := col.FindOne(ctx, bson.D{}, opts).Decode(&round)
	if err == mongo.ErrNoDocuments {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return round.VoteId, nil
}

// InsertVoteIfNotExists inserts a vote; no-ops on (politicianId, roundId) duplicate.
func (database *DB) InsertVoteIfNotExists(ctx context.Context, vote *models.Vote) error {
	col := database.database.Collection(ColVotes)
	_, err := col.UpdateOne(
		ctx,
		bson.D{{Key: "politicianId", Value: vote.PoliticianId}, {Key: "roundId", Value: vote.RoundId}},
		bson.D{{Key: "$setOnInsert", Value: vote}},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}

// GetRoundsWithoutDescription returns rounds where description is empty (need enrichment).
func (database *DB) GetRoundsWithoutDescription(ctx context.Context, limit int) ([]models.VotingRound, error) {
	col := database.database.Collection(ColVotingRounds)
	filter := bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "description", Value: ""}},
		bson.D{{Key: "description", Value: bson.D{{Key: "$exists", Value: false}}}},
	}}}
	cursor, err := col.Find(ctx, filter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	var rounds []models.VotingRound
	return rounds, cursor.All(ctx, &rounds)
}

// UpdateRoundDescription sets the title and description for a round.
func (database *DB) UpdateRoundDescription(ctx context.Context, id, title, description string) error {
	col := database.database.Collection(ColVotingRounds)
	_, err := col.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "title", Value: title},
			{Key: "name", Value: title},
			{Key: "description", Value: description},
		}}},
	)
	return err
}

// ReplaceActivePoliticians deactivates all politicians, then reactivates those in the given name list.
func (database *DB) ReplaceActivePoliticians(ctx context.Context, activeNames []string) error {
	col := database.database.Collection(ColPoliticians)

	// Deactivate all
	if _, err := col.UpdateMany(ctx, bson.D{}, bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: false}}}}); err != nil {
		return fmt.Errorf("deactivate all: %w", err)
	}

	if len(activeNames) == 0 {
		return nil
	}

	// Reactivate matching names
	_, err := col.UpdateMany(
		ctx,
		bson.D{{Key: "name", Value: bson.D{{Key: "$in", Value: activeNames}}}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: true}}}},
	)
	return err
}
