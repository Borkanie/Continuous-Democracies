package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/borkanie/brand-new-day-api/internal/db"
	"github.com/borkanie/brand-new-day-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// PoliticianRepository is the read-only access surface for the politicians
// collection.
type PoliticianRepository interface {
	ListPoliticians(requestContext context.Context) ([]models.Politician, error)
	GetPoliticianByID(requestContext context.Context, politicianID string) (*models.Politician, error)
	GetPoliticiansByIDs(requestContext context.Context, politicianIDs []string) (map[string]models.Politician, error)
}

// PoliticianRepositoryMongo is the MongoDB-backed PoliticianRepository
// implementation.
type PoliticianRepositoryMongo struct {
	politiciansCollection *mongo.Collection
}

var _ PoliticianRepository = (*PoliticianRepositoryMongo)(nil)

// NewPoliticianRepository builds a PoliticianRepositoryMongo bound to the
// given database's politicians collection.
func NewPoliticianRepository(database *mongo.Database) *PoliticianRepositoryMongo {
	return &PoliticianRepositoryMongo{
		politiciansCollection: database.Collection(db.CollectionPoliticians),
	}
}

// ListPoliticians returns every politician. Never returns a nil slice.
func (politicianRepository *PoliticianRepositoryMongo) ListPoliticians(requestContext context.Context) ([]models.Politician, error) {
	cursor, err := politicianRepository.politiciansCollection.Find(requestContext, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("list politicians: %w", err)
	}

	politicians := make([]models.Politician, 0)
	if err := cursor.All(requestContext, &politicians); err != nil {
		return nil, fmt.Errorf("decode politicians: %w", err)
	}
	return politicians, nil
}

// GetPoliticianByID returns the politician with the given id, or ErrNotFound
// if none exists.
func (politicianRepository *PoliticianRepositoryMongo) GetPoliticianByID(requestContext context.Context, politicianID string) (*models.Politician, error) {
	var politician models.Politician
	err := politicianRepository.politiciansCollection.FindOne(requestContext, bson.D{{Key: "_id", Value: politicianID}}).Decode(&politician)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get politician by id: %w", err)
	}
	return &politician, nil
}

// GetPoliticiansByIDs returns every requested politician keyed by id. Missing
// ids are simply absent from the result map; this is a bulk lookup helper,
// not a per-id existence check.
func (politicianRepository *PoliticianRepositoryMongo) GetPoliticiansByIDs(requestContext context.Context, politicianIDs []string) (map[string]models.Politician, error) {
	politiciansByID := make(map[string]models.Politician, len(politicianIDs))
	if len(politicianIDs) == 0 {
		return politiciansByID, nil
	}

	cursor, err := politicianRepository.politiciansCollection.Find(requestContext, bson.D{
		{Key: "_id", Value: bson.D{{Key: "$in", Value: politicianIDs}}},
	})
	if err != nil {
		return nil, fmt.Errorf("get politicians by ids: %w", err)
	}

	var politicians []models.Politician
	if err := cursor.All(requestContext, &politicians); err != nil {
		return nil, fmt.Errorf("decode politicians: %w", err)
	}

	for _, politician := range politicians {
		politiciansByID[politician.ID] = politician
	}
	return politiciansByID, nil
}
