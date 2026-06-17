package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/borkanie/parliament-api/internal/db"
	"github.com/borkanie/parliament-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type VotingRepo struct {
	roundsCol *mongo.Collection
	votesCol  *mongo.Collection
	database  *mongo.Database
}

func NewVotingRepo(database *mongo.Database) *VotingRepo {
	return &VotingRepo{
		roundsCol: database.Collection(db.ColVotingRounds),
		votesCol:  database.Collection(db.ColVotes),
		database:  database,
	}
}

type RoundFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Keywords  []string
	Limit     int
}

func (repo *VotingRepo) GetAllRounds(ctx context.Context, roundFilter RoundFilter) ([]models.VotingRound, error) {
	filter := bson.D{}

	if roundFilter.StartDate != nil || roundFilter.EndDate != nil {
		dateFilter := bson.D{}
		if roundFilter.StartDate != nil {
			dateFilter = append(dateFilter, bson.E{Key: "$gte", Value: *roundFilter.StartDate})
		}
		if roundFilter.EndDate != nil {
			dateFilter = append(dateFilter, bson.E{Key: "$lte", Value: *roundFilter.EndDate})
		}
		filter = append(filter, bson.E{Key: "voteDate", Value: dateFilter})
	}

	if len(roundFilter.Keywords) > 0 {
		// Use MongoDB text search for the first keyword; fall back to regex for multiple
		if len(roundFilter.Keywords) == 1 {
			filter = append(filter, bson.E{Key: "$text", Value: bson.D{{Key: "$search", Value: roundFilter.Keywords[0]}}})
		} else {
			regexOrs := make(bson.A, 0, len(roundFilter.Keywords))
			for _, keyword := range roundFilter.Keywords {
				regexOrs = append(regexOrs, bson.D{{Key: "$or", Value: bson.A{
					bson.D{{Key: "title", Value: bson.D{{Key: "$regex", Value: keyword}, {Key: "$options", Value: "i"}}}},
					bson.D{{Key: "description", Value: bson.D{{Key: "$regex", Value: keyword}, {Key: "$options", Value: "i"}}}},
				}}})
			}
			filter = append(filter, bson.E{Key: "$and", Value: regexOrs})
		}
	}

	limit := 100
	if roundFilter.Limit > 0 {
		limit = roundFilter.Limit
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "voteDate", Value: -1}})

	cursor, err := repo.roundsCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("rounds find: %w", err)
	}

	var rounds []models.VotingRound
	if err := cursor.All(ctx, &rounds); err != nil {
		return nil, fmt.Errorf("rounds decode: %w", err)
	}
	if rounds == nil {
		rounds = []models.VotingRound{}
	}
	return rounds, nil
}

func (repo *VotingRepo) GetRoundByVoteId(ctx context.Context, voteId int) (*models.VotingRound, error) {
	var round models.VotingRound
	err := repo.roundsCol.FindOne(ctx, bson.D{{Key: "voteId", Value: voteId}}).Decode(&round)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("round findByVoteId: %w", err)
	}
	return &round, nil
}

type VoteFilter struct {
	// Either RoundVoteId (int) or RoundID (UUID string) must be set
	RoundVoteId *int
	RoundID     string
	PartyId     string
	PartyAcronim string
}

func (repo *VotingRepo) GetVotes(ctx context.Context, voteFilter VoteFilter) ([]models.VoteResult, error) {
	// Resolve roundId from voteId if needed
	roundId := voteFilter.RoundID
	if voteFilter.RoundVoteId != nil {
		round, err := repo.GetRoundByVoteId(ctx, *voteFilter.RoundVoteId)
		if err != nil {
			return nil, err
		}
		if round == nil {
			return []models.VoteResult{}, nil
		}
		roundId = round.ID
	}

	pipeline := bson.A{
		// Match votes for this round
		bson.D{{Key: "$match", Value: bson.D{{Key: "roundId", Value: roundId}}}},

		// Join politician
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: db.ColPoliticians},
			{Key: "localField", Value: "politicianId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "politician"},
		}}},
		bson.D{{Key: "$unwind", Value: "$politician"}},

		// Join party into politician
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: db.ColParties},
			{Key: "localField", Value: "politician.partyId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "politician.party"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$politician.party"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},

		// Join round
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: db.ColVotingRounds},
			{Key: "localField", Value: "roundId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "round"},
		}}},
		bson.D{{Key: "$unwind", Value: "$round"}},

		// Add computed name field
		bson.D{{Key: "$addFields", Value: bson.D{{Key: "name", Value: "$politician.name"}}}},
	}

	// Filter by party after joins
	if voteFilter.PartyId != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "politician.partyId", Value: voteFilter.PartyId}}}})
	}
	if voteFilter.PartyAcronim != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
			{Key: "politician.party.acronym", Value: bson.D{
				{Key: "$regex", Value: voteFilter.PartyAcronim},
				{Key: "$options", Value: "i"},
			}},
		}}})
	}

	cursor, err := repo.votesCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("votes aggregate: %w", err)
	}

	var results []models.VoteResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("votes decode: %w", err)
	}
	if results == nil {
		results = []models.VoteResult{}
	}
	return results, nil
}
