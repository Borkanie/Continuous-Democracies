package repository

import (
	"context"
	"fmt"

	"github.com/borkanie/brand-new-day-api/internal/db"
	"github.com/borkanie/brand-new-day-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// LawBucketRepository is the read-only access surface for the law_buckets
// collection.
type LawBucketRepository interface {
	ListLawBuckets(requestContext context.Context) ([]models.LawBucket, error)
	GetLawBucketByID(requestContext context.Context, lawBucketID int) (*models.LawBucket, error)
}

// LawBucketRepositoryMongo is the MongoDB-backed LawBucketRepository
// implementation.
type LawBucketRepositoryMongo struct {
	lawBucketsCollection db.MongoCollection
}

var _ LawBucketRepository = (*LawBucketRepositoryMongo)(nil)

// NewLawBucketRepository builds a LawBucketRepositoryMongo bound to the
// given database's law_buckets collection.
func NewLawBucketRepository(database *mongo.Database) *LawBucketRepositoryMongo {
	return &LawBucketRepositoryMongo{
		lawBucketsCollection: db.GetLoggingCollection(database, db.CollectionLawBuckets),
	}
}

// ListLawBuckets returns every law bucket, each with .Normatives populated
// with the current (highest-version) normative for every distinct normative
// id belonging to that bucket. Never returns a nil slice.
func (lawBucketRepository *LawBucketRepositoryMongo) ListLawBuckets(requestContext context.Context) ([]models.LawBucket, error) {
	pipeline := bson.A{currentNormativesLookupStage()}

	cursor, err := lawBucketRepository.lawBucketsCollection.Aggregate(requestContext, pipeline, options.Aggregate())
	if err != nil {
		return nil, fmt.Errorf("aggregate law buckets: %w", err)
	}

	lawBuckets := make([]models.LawBucket, 0)
	if err := cursor.All(requestContext, &lawBuckets); err != nil {
		return nil, fmt.Errorf("decode law buckets: %w", err)
	}
	return lawBuckets, nil
}

// GetLawBucketByID returns the law bucket with the given id, with
// .Normatives populated the same way as ListLawBuckets, or ErrNotFound if
// none exists.
func (lawBucketRepository *LawBucketRepositoryMongo) GetLawBucketByID(requestContext context.Context, lawBucketID int) (*models.LawBucket, error) {
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: lawBucketID}}}},
		currentNormativesLookupStage(),
		bson.D{{Key: "$limit", Value: 1}},
	}

	cursor, err := lawBucketRepository.lawBucketsCollection.Aggregate(requestContext, pipeline, options.Aggregate())
	if err != nil {
		return nil, fmt.Errorf("aggregate law bucket by id: %w", err)
	}
	defer cursor.Close(requestContext)

	if !cursor.Next(requestContext) {
		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("aggregate law bucket by id cursor: %w", err)
		}
		return nil, ErrNotFound
	}

	var lawBucket models.LawBucket
	if err := cursor.Decode(&lawBucket); err != nil {
		return nil, fmt.Errorf("decode law bucket by id: %w", err)
	}
	return &lawBucket, nil
}

// currentNormativesLookupStage builds the $lookup that embeds, on every
// law_buckets document flowing through the pipeline, the current
// (highest-version) normative for each distinct normative id belonging to
// that bucket, ordered by label. It matches the bucket via an $expr on the
// joined normative's lawBucketId (the "let"/pipeline form of $lookup),
// because the reduction to "current version" needs extra stages
// (sort+group+replaceRoot) beyond a plain equality join.
func currentNormativesLookupStage() bson.D {
	subPipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$expr", Value: bson.D{{Key: "$eq", Value: bson.A{"$lawBucketId", "$$bucketId"}}}},
		}}},
	}
	subPipeline = append(subPipeline, reduceToCurrentNormativeStages()...)

	return bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: db.CollectionNormatives},
		{Key: "let", Value: bson.D{{Key: "bucketId", Value: "$_id"}}},
		{Key: "pipeline", Value: subPipeline},
		{Key: "as", Value: "normatives"},
	}}}
}
