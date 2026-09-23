package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/borkanie/brand-new-day-api/internal/db"
	"github.com/borkanie/brand-new-day-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NormativeRepository is the read-only access surface for the normatives
// collection.
type NormativeRepository interface {
	GetNormativeByIDAndVersion(requestContext context.Context, normativeID string, normativeVersion int) (*models.Normative, error)
	ListCurrentNormativesByLawBucketID(requestContext context.Context, lawBucketID int) ([]models.Normative, error)
}

// NormativeRepositoryMongo is the MongoDB-backed NormativeRepository
// implementation.
type NormativeRepositoryMongo struct {
	normativesCollection db.MongoCollection
}

var _ NormativeRepository = (*NormativeRepositoryMongo)(nil)

// NewNormativeRepository builds a NormativeRepositoryMongo bound to the given
// database's normatives collection.
func NewNormativeRepository(database *mongo.Database) *NormativeRepositoryMongo {
	return &NormativeRepositoryMongo{
		normativesCollection: db.GetLoggingCollection(database, db.CollectionNormatives),
	}
}

// GetNormativeByIDAndVersion returns the exact normative version identified
// by the (id, version) FK pair, or ErrNotFound if none exists.
func (normativeRepository *NormativeRepositoryMongo) GetNormativeByIDAndVersion(requestContext context.Context, normativeID string, normativeVersion int) (*models.Normative, error) {
	var normative models.Normative
	filter := bson.D{
		{Key: "_id.id", Value: normativeID},
		{Key: "_id.version", Value: normativeVersion},
	}
	err := normativeRepository.normativesCollection.FindOne(requestContext, filter).Decode(&normative)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get normative by id and version: %w", err)
	}
	return &normative, nil
}

// ListCurrentNormativesByLawBucketID returns the current (highest-version)
// normative for each distinct normative id belonging to lawBucketID, ordered
// deterministically by label. Never returns a nil slice.
func (normativeRepository *NormativeRepositoryMongo) ListCurrentNormativesByLawBucketID(requestContext context.Context, lawBucketID int) ([]models.Normative, error) {
	pipeline := currentNormativesPipeline(lawBucketID)

	cursor, err := normativeRepository.normativesCollection.Aggregate(requestContext, pipeline, options.Aggregate())
	if err != nil {
		return nil, fmt.Errorf("aggregate current normatives: %w", err)
	}

	normatives := make([]models.Normative, 0)
	if err := cursor.All(requestContext, &normatives); err != nil {
		return nil, fmt.Errorf("decode current normatives: %w", err)
	}
	return normatives, nil
}

// currentNormativesPipeline builds the aggregation stages that reduce every
// version of every normative belonging to lawBucketID down to just the
// current (highest-version) document per distinct normative id, sorted by
// label.
func currentNormativesPipeline(lawBucketID int) bson.A {
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "lawBucketId", Value: lawBucketID}}}},
	}
	return append(pipeline, reduceToCurrentNormativeStages()...)
}

// reduceToCurrentNormativeStages returns the $sort/$group/$replaceRoot/$sort
// stages that pick the highest-version document per distinct normative id
// and order the result by label. It assumes an earlier stage has already
// narrowed the pipeline to the normatives of a single law bucket (either a
// plain $match, as in currentNormativesPipeline, or an $expr match inside a
// $lookup sub-pipeline, as in lawbucket.go).
func reduceToCurrentNormativeStages() bson.A {
	return bson.A{
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id.version", Value: -1}}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id.id"},
			{Key: "current", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
		}}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$current"}}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "label", Value: 1}}}},
	}
}
