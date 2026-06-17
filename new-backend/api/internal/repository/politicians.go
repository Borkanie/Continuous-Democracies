package repository

import (
	"context"
	"fmt"

	"github.com/borkanie/parliament-api/internal/db"
	"github.com/borkanie/parliament-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PoliticiansRepo struct {
	col      *mongo.Collection
	database *mongo.Database
}

func NewPoliticiansRepo(database *mongo.Database) *PoliticiansRepo {
	return &PoliticiansRepo{
		col:      database.Collection(db.ColPoliticians),
		database: database,
	}
}

type PoliticianFilter struct {
	PartyAcronym string
	PartyName    string
	IsActive     *bool
	Location     *int
	Gender       *int
	Limit        int
}

func (repo *PoliticiansRepo) GetAll(ctx context.Context, filter PoliticianFilter) ([]models.PoliticianWithParty, error) {
	pipeline := bson.A{}

	// Join parties first so we can filter by party name/acronym
	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: db.ColParties},
		{Key: "localField", Value: "partyId"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "party"},
	}}})
	pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$party"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}})

	// Build match stage
	matchStage := bson.D{}
	if filter.IsActive != nil {
		matchStage = append(matchStage, bson.E{Key: "active", Value: *filter.IsActive})
	}
	if filter.Location != nil {
		matchStage = append(matchStage, bson.E{Key: "workLocation", Value: *filter.Location})
	}
	if filter.Gender != nil {
		matchStage = append(matchStage, bson.E{Key: "gender", Value: *filter.Gender})
	}
	if filter.PartyAcronym != "" {
		matchStage = append(matchStage, bson.E{Key: "party.acronym", Value: bson.D{
			{Key: "$regex", Value: filter.PartyAcronym},
			{Key: "$options", Value: "i"},
		}})
	}
	if filter.PartyName != "" {
		matchStage = append(matchStage, bson.E{Key: "party.name", Value: bson.D{
			{Key: "$regex", Value: filter.PartyName},
			{Key: "$options", Value: "i"},
		}})
	}

	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchStage}})
	}

	if filter.Limit > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: filter.Limit}})
	}

	cursor, err := repo.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("politician aggregate: %w", err)
	}

	var results []models.PoliticianWithParty
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("politician decode: %w", err)
	}
	if results == nil {
		results = []models.PoliticianWithParty{}
	}
	return results, nil
}

func (repo *PoliticiansRepo) GetByID(ctx context.Context, id string) (*models.PoliticianWithParty, error) {
	return repo.getOneWithParty(ctx, bson.D{{Key: "_id", Value: id}})
}

func (repo *PoliticiansRepo) GetByName(ctx context.Context, name string) (*models.PoliticianWithParty, error) {
	return repo.getOneWithParty(ctx, bson.D{{Key: "name", Value: bson.D{
		{Key: "$regex", Value: name},
		{Key: "$options", Value: "i"},
	}}})
}

func (repo *PoliticiansRepo) getOneWithParty(ctx context.Context, matchFilter bson.D) (*models.PoliticianWithParty, error) {
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: matchFilter}},
		bson.D{{Key: "$limit", Value: 1}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: db.ColParties},
			{Key: "localField", Value: "partyId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "party"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$party"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := repo.col.Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return nil, fmt.Errorf("politician getOne aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, nil
	}

	var politician models.PoliticianWithParty
	if err := cursor.Decode(&politician); err != nil {
		return nil, fmt.Errorf("politician getOne decode: %w", err)
	}
	return &politician, nil
}
