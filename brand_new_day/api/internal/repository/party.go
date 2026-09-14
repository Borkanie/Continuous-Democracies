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

// PartyRepository is the read-only access surface for the parties
// collection.
type PartyRepository interface {
	ListParties(requestContext context.Context) ([]models.Party, error)
	GetPartyByID(requestContext context.Context, partyID string) (*models.Party, error)
	GetPartiesByIDs(requestContext context.Context, partyIDs []string) (map[string]models.Party, error)
}

// PartyRepositoryMongo is the MongoDB-backed PartyRepository implementation.
type PartyRepositoryMongo struct {
	partiesCollection *mongo.Collection
}

var _ PartyRepository = (*PartyRepositoryMongo)(nil)

// NewPartyRepository builds a PartyRepositoryMongo bound to the given
// database's parties collection.
func NewPartyRepository(database *mongo.Database) *PartyRepositoryMongo {
	return &PartyRepositoryMongo{
		partiesCollection: database.Collection(db.CollectionParties),
	}
}

// ListParties returns every party. Never returns a nil slice.
func (partyRepository *PartyRepositoryMongo) ListParties(requestContext context.Context) ([]models.Party, error) {
	cursor, err := partyRepository.partiesCollection.Find(requestContext, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("list parties: %w", err)
	}

	parties := make([]models.Party, 0)
	if err := cursor.All(requestContext, &parties); err != nil {
		return nil, fmt.Errorf("decode parties: %w", err)
	}
	return parties, nil
}

// GetPartyByID returns the party with the given id, or ErrNotFound if none
// exists.
func (partyRepository *PartyRepositoryMongo) GetPartyByID(requestContext context.Context, partyID string) (*models.Party, error) {
	var party models.Party
	err := partyRepository.partiesCollection.FindOne(requestContext, bson.D{{Key: "_id", Value: partyID}}).Decode(&party)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get party by id: %w", err)
	}
	return &party, nil
}

// GetPartiesByIDs returns every requested party keyed by id. Missing ids are
// simply absent from the result map; this is a bulk lookup helper, not a
// per-id existence check.
func (partyRepository *PartyRepositoryMongo) GetPartiesByIDs(requestContext context.Context, partyIDs []string) (map[string]models.Party, error) {
	partiesByID := make(map[string]models.Party, len(partyIDs))
	if len(partyIDs) == 0 {
		return partiesByID, nil
	}

	cursor, err := partyRepository.partiesCollection.Find(requestContext, bson.D{
		{Key: "_id", Value: bson.D{{Key: "$in", Value: partyIDs}}},
	})
	if err != nil {
		return nil, fmt.Errorf("get parties by ids: %w", err)
	}

	var parties []models.Party
	if err := cursor.All(requestContext, &parties); err != nil {
		return nil, fmt.Errorf("decode parties: %w", err)
	}

	for _, party := range parties {
		partiesByID[party.ID] = party
	}
	return partiesByID, nil
}
