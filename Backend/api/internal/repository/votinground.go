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

// VotingRoundRepository is the read-only access surface for the
// voting_rounds collection. It owns every vote-array query, including the
// "all votes by a given politician" aggregation, because Vote is embedded in
// VotingRound rather than being its own collection.
type VotingRoundRepository interface {
	ListVotingRounds(requestContext context.Context) ([]models.VotingRound, error)
	GetVotingRoundByID(requestContext context.Context, votingRoundID int) (*models.VotingRound, error)
	ListVotingRoundsByLawBucketID(requestContext context.Context, lawBucketID int) ([]models.VotingRound, error)
	ListVoteRecordsByPoliticianID(requestContext context.Context, politicianID string) ([]models.PoliticianVoteRecord, error)
}

// VotingRoundRepositoryMongo is the MongoDB-backed VotingRoundRepository
// implementation.
type VotingRoundRepositoryMongo struct {
	votingRoundsCollection db.MongoCollection
	normativesCollection   db.MongoCollection
}

var _ VotingRoundRepository = (*VotingRoundRepositoryMongo)(nil)

// NewVotingRoundRepository builds a VotingRoundRepositoryMongo bound to the
// given database's voting_rounds collection (and, for
// ListVotingRoundsByLawBucketID's two-step join, the normatives collection).
func NewVotingRoundRepository(database *mongo.Database) *VotingRoundRepositoryMongo {
	return &VotingRoundRepositoryMongo{
		votingRoundsCollection: db.GetLoggingCollection(database, db.CollectionVotingRounds),
		normativesCollection:   db.GetLoggingCollection(database, db.CollectionNormatives),
	}
}

// ListVotingRounds returns every voting round. Never returns a nil slice.
func (votingRoundRepository *VotingRoundRepositoryMongo) ListVotingRounds(requestContext context.Context) ([]models.VotingRound, error) {
	cursor, err := votingRoundRepository.votingRoundsCollection.Find(requestContext, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("list voting rounds: %w", err)
	}

	votingRounds := make([]models.VotingRound, 0)
	if err := cursor.All(requestContext, &votingRounds); err != nil {
		return nil, fmt.Errorf("decode voting rounds: %w", err)
	}
	return votingRounds, nil
}

// GetVotingRoundByID returns the voting round with the given id, or
// ErrNotFound if none exists.
func (votingRoundRepository *VotingRoundRepositoryMongo) GetVotingRoundByID(requestContext context.Context, votingRoundID int) (*models.VotingRound, error) {
	var votingRound models.VotingRound
	err := votingRoundRepository.votingRoundsCollection.FindOne(requestContext, bson.D{{Key: "_id", Value: votingRoundID}}).Decode(&votingRound)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get voting round by id: %w", err)
	}
	return &votingRound, nil
}

// ListVotingRoundsByLawBucketID returns every voting round whose
// normativeId belongs to any normative of the given law bucket. Implemented
// as a two-step lookup (distinct normative ids for the bucket, then a
// voting_rounds find with $in) rather than a single aggregation, per the
// contract's "either is fine" allowance. Never returns a nil slice.
func (votingRoundRepository *VotingRoundRepositoryMongo) ListVotingRoundsByLawBucketID(requestContext context.Context, lawBucketID int) ([]models.VotingRound, error) {
	distinctResult := votingRoundRepository.normativesCollection.Distinct(
		requestContext,
		"_id.id",
		bson.D{{Key: "lawBucketId", Value: lawBucketID}},
		options.Distinct(),
	)

	var normativeIDs []string
	if err := distinctResult.Decode(&normativeIDs); err != nil {
		return nil, fmt.Errorf("distinct normative ids for law bucket: %w", err)
	}
	if len(normativeIDs) == 0 {
		return []models.VotingRound{}, nil
	}

	cursor, err := votingRoundRepository.votingRoundsCollection.Find(requestContext, bson.D{
		{Key: "normativeId", Value: bson.D{{Key: "$in", Value: normativeIDs}}},
	})
	if err != nil {
		return nil, fmt.Errorf("find voting rounds by normative ids: %w", err)
	}

	votingRounds := make([]models.VotingRound, 0)
	if err := cursor.All(requestContext, &votingRounds); err != nil {
		return nil, fmt.Errorf("decode voting rounds by law bucket: %w", err)
	}
	return votingRounds, nil
}

// ListVoteRecordsByPoliticianID is the headline aggregation: every vote cast
// by politicianID, across every voting round, enriched with the round, the
// exact normative version voted on, and the law bucket it resolves to.
// Sorted by voteDate descending. Never returns a nil slice.
func (votingRoundRepository *VotingRoundRepositoryMongo) ListVoteRecordsByPoliticianID(requestContext context.Context, politicianID string) ([]models.PoliticianVoteRecord, error) {
	politicianVoteFilter := bson.D{{Key: "votes.politicianId", Value: politicianID}}

	pipeline := bson.A{
		// Uses the multikey index on votes.politicianId.
		bson.D{{Key: "$match", Value: politicianVoteFilter}},
		bson.D{{Key: "$unwind", Value: "$votes"}},
		bson.D{{Key: "$match", Value: politicianVoteFilter}},

		// Join the exact normative version this round voted on. A compound
		// _id ({id, version}) join requires the let/pipeline form of
		// $lookup; localField/foreignField can't express the pair match.
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: db.CollectionNormatives},
			{Key: "let", Value: bson.D{
				{Key: "normativeId", Value: "$normativeId"},
				{Key: "normativeVersion", Value: "$normativeVersion"},
			}},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{{Key: "$and", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$_id.id", "$$normativeId"}}},
						bson.D{{Key: "$eq", Value: bson.A{"$_id.version", "$$normativeVersion"}}},
					}}}},
				}}},
			}},
			{Key: "as", Value: "normative"},
		}}},
		// preserveNullAndEmptyArrays keeps the vote in the result even if its
		// normative reference is dangling — PoliticianVoteRecord.Normative is
		// a nil-able pointer for exactly this case. Without it, a broken
		// reference would silently drop the vote from the politician's
		// history instead of surfacing it with the law info missing.
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$normative"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},

		// Resolve the normative to its parent law bucket.
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: db.CollectionLawBuckets},
			{Key: "localField", Value: "normative.lawBucketId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "lawBucket"},
		}}},
		// Same rationale as the normative $unwind above.
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$lawBucket"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},

		// Sort while the round's own voteDate is still a top-level field
		// (before $project reshapes the document).
		bson.D{{Key: "$sort", Value: bson.D{{Key: "voteDate", Value: -1}}}},

		// Project into the PoliticianVoteRecord shape. votingRound.votes is
		// forced to an empty array so a politician's vote history isn't
		// ballooned by every other politician's vote on the same round.
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "value", Value: "$votes.value"},
			{Key: "votingRound", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "title", Value: "$title"},
				{Key: "description", Value: "$description"},
				{Key: "voteDate", Value: "$voteDate"},
				{Key: "normativeId", Value: "$normativeId"},
				{Key: "normativeVersion", Value: "$normativeVersion"},
				{Key: "chamber", Value: "$chamber"},
				{Key: "votes", Value: bson.D{{Key: "$literal", Value: bson.A{}}}},
			}},
			{Key: "normative", Value: "$normative"},
			{Key: "lawBucket", Value: "$lawBucket"},
		}}},
	}

	cursor, err := votingRoundRepository.votingRoundsCollection.Aggregate(requestContext, pipeline, options.Aggregate())
	if err != nil {
		return nil, fmt.Errorf("aggregate vote records by politician: %w", err)
	}

	voteRecords := make([]models.PoliticianVoteRecord, 0)
	if err := cursor.All(requestContext, &voteRecords); err != nil {
		return nil, fmt.Errorf("decode vote records by politician: %w", err)
	}
	return voteRecords, nil
}
