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

type PartyRepo struct {
	col *mongo.Collection
}

func NewPartyRepo(database *mongo.Database) *PartyRepo {
	return &PartyRepo{col: database.Collection(db.ColParties)}
}

func (repo *PartyRepo) All(ctx context.Context, active *bool, limit int) ([]models.Party, error) {
	filter := bson.D{}
	if active != nil {
		filter = append(filter, bson.E{Key: "active", Value: *active})
	}

	opts := options.Find().SetLimit(int64(limit))
	cursor, err := repo.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("party find: %w", err)
	}

	var parties []models.Party
	if err := cursor.All(ctx, &parties); err != nil {
		return nil, fmt.Errorf("party decode: %w", err)
	}
	if parties == nil {
		parties = []models.Party{}
	}
	return parties, nil
}

func (repo *PartyRepo) GetByID(ctx context.Context, id string) (*models.Party, error) {
	var party models.Party
	err := repo.col.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&party)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("party findById: %w", err)
	}
	return &party, nil
}

func (repo *PartyRepo) Query(ctx context.Context, name, acronym string) (*models.Party, error) {
	filter := bson.D{}
	if name != "" {
		filter = append(filter, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: name}, {Key: "$options", Value: "i"}}})
	}
	if acronym != "" {
		filter = append(filter, bson.E{Key: "acronym", Value: bson.D{{Key: "$regex", Value: acronym}, {Key: "$options", Value: "i"}}})
	}

	var party models.Party
	err := repo.col.FindOne(ctx, filter).Decode(&party)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("party query: %w", err)
	}
	return &party, nil
}
