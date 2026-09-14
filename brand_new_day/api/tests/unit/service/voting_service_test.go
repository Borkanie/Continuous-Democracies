package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/borkanie/brand-new-day-api/internal/models"
	"github.com/borkanie/brand-new-day-api/internal/repository"
	"github.com/borkanie/brand-new-day-api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestVotingServiceListVotingRounds(testRunner *testing.T) {
	testRunner.Run("returns voting rounds with empty votes slice, not nil", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")

		votingRoundModelsToReturn := []models.VotingRound{
			{
				ID:               1,
				Title:            "Round 1",
				Description:      "First round",
				VoteDate:         voteDate,
				NormativeID:      "norm-1",
				NormativeVersion: 1,
				Votes: []models.Vote{
					{PoliticianID: "pol-1", PartyID: "party-1", Value: "Yes"},
					{PoliticianID: "pol-2", PartyID: "party-1", Value: "No"},
					{PoliticianID: "pol-3", PartyID: "party-2", Value: "Abstain"},
				},
			},
		}

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			listVotingRoundsResult: votingRoundModelsToReturn,
		}
		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{}
		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		votingServiceInstance := service.NewVotingService(
			fakeVotingRoundRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeLawBucketRepositoryInstance,
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
		)

		votingRoundDTOsResult, errResult := votingServiceInstance.ListVotingRounds(testContext)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, votingRoundDTOsResult, 1)
		// Votes should be an empty slice, not nil
		require.NotNil(testRunner, votingRoundDTOsResult[0].Votes)
		require.Len(testRunner, votingRoundDTOsResult[0].Votes, 0)
		require.Equal(testRunner, 1, votingRoundDTOsResult[0].Id)
		require.Equal(testRunner, "Round 1", votingRoundDTOsResult[0].Title)
	})

	testRunner.Run("returns empty slice when no voting rounds exist", func(testRunner *testing.T) {
		testContext := context.Background()

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			listVotingRoundsResult: []models.VotingRound{},
		}
		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{}
		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		votingServiceInstance := service.NewVotingService(
			fakeVotingRoundRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeLawBucketRepositoryInstance,
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
		)

		votingRoundDTOsResult, errResult := votingServiceInstance.ListVotingRounds(testContext)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, votingRoundDTOsResult, 0)
	})
}

func TestVotingServiceGetVotingRoundByID(testRunner *testing.T) {
	testRunner.Run("returns voting round detail with resolved normative and law bucket", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")
		effectiveDate, _ := time.Parse("2006-01-02", "2024-01-10")
		initiationDate, _ := time.Parse("2006-01-02", "2024-01-01")

		votingRoundModelToReturn := &models.VotingRound{
			ID:               1,
			Title:            "Round 1",
			Description:      "First round",
			VoteDate:         voteDate,
			NormativeID:      "norm-1",
			NormativeVersion: 1,
			Votes: []models.Vote{
				{PoliticianID: "pol-1", PartyID: "party-1", Value: "Yes"},
			},
		}

		normativeModelToReturn := &models.Normative{
			Key:           models.NormativeKey{ID: "norm-1", Version: 1},
			LawBucketID:   42,
			Type:          "article",
			Label:         "Article 5",
			Text:          "Some text",
			EffectiveDate: effectiveDate,
		}

		lawBucketModelToReturn := &models.LawBucket{
			ID:             42,
			Version:        1,
			PLNumber:       "PL 123/2024",
			Title:          "Test Law",
			Description:    "Test description",
			InitiationDate: initiationDate,
			Status:         "pending",
			Normatives:     []models.Normative{},
		}

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			getVotingRoundByIDFunc: func(ctx context.Context, votingRoundID int) (*models.VotingRound, error) {
				if votingRoundID == 1 {
					return votingRoundModelToReturn, nil
				}
				return nil, repository.ErrNotFound
			},
		}

		fakeNormativeRepositoryInstance := &fakeNormativeRepository{
			getNormativeByIDAndVersionFunc: func(ctx context.Context, normativeID string, normativeVersion int) (*models.Normative, error) {
				if normativeID == "norm-1" && normativeVersion == 1 {
					return normativeModelToReturn, nil
				}
				return nil, repository.ErrNotFound
			},
		}

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{
			getLawBucketByIDFunc: func(ctx context.Context, lawBucketID int) (*models.LawBucket, error) {
				if lawBucketID == 42 {
					return lawBucketModelToReturn, nil
				}
				return nil, repository.ErrNotFound
			},
		}

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		votingServiceInstance := service.NewVotingService(
			fakeVotingRoundRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeLawBucketRepositoryInstance,
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
		)

		votingRoundDetailDTOResult, errResult := votingServiceInstance.GetVotingRoundByID(testContext, 1)

		require.NoError(testRunner, errResult)
		require.NotNil(testRunner, votingRoundDetailDTOResult)
		require.Equal(testRunner, 1, votingRoundDetailDTOResult.Id)
		require.Equal(testRunner, "norm-1", votingRoundDetailDTOResult.Normative.Id)
		require.Equal(testRunner, 1, votingRoundDetailDTOResult.Normative.Version)
		require.Equal(testRunner, "Article 5", votingRoundDetailDTOResult.Normative.Label)
		require.Equal(testRunner, 42, votingRoundDetailDTOResult.LawBucket.Id)
		require.Equal(testRunner, "PL 123/2024", votingRoundDetailDTOResult.LawBucket.PlNumber)
	})

	testRunner.Run("returns partial data when normative is missing", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")

		votingRoundModelToReturn := &models.VotingRound{
			ID:               1,
			Title:            "Round 1",
			VoteDate:         voteDate,
			NormativeID:      "norm-missing",
			NormativeVersion: 1,
			Votes:            []models.Vote{},
		}

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			getVotingRoundByIDFunc: func(ctx context.Context, votingRoundID int) (*models.VotingRound, error) {
				return votingRoundModelToReturn, nil
			},
		}

		fakeNormativeRepositoryInstance := &fakeNormativeRepository{
			getNormativeByIDAndVersionFunc: func(ctx context.Context, normativeID string, normativeVersion int) (*models.Normative, error) {
				return nil, repository.ErrNotFound
			},
		}

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{}
		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		votingServiceInstance := service.NewVotingService(
			fakeVotingRoundRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeLawBucketRepositoryInstance,
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
		)

		votingRoundDetailDTOResult, errResult := votingServiceInstance.GetVotingRoundByID(testContext, 1)

		require.NoError(testRunner, errResult)
		require.NotNil(testRunner, votingRoundDetailDTOResult)
		require.Equal(testRunner, 1, votingRoundDetailDTOResult.Id)
		// Normative should be zero-valued
		require.Equal(testRunner, "", votingRoundDetailDTOResult.Normative.Id)
		// Law bucket should be zero-valued because normative wasn't resolved
		require.Equal(testRunner, 0, votingRoundDetailDTOResult.LawBucket.Id)
	})

	testRunner.Run("returns partial data when law bucket is missing", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")
		effectiveDate, _ := time.Parse("2006-01-02", "2024-01-10")

		votingRoundModelToReturn := &models.VotingRound{
			ID:               1,
			Title:            "Round 1",
			VoteDate:         voteDate,
			NormativeID:      "norm-1",
			NormativeVersion: 1,
			Votes:            []models.Vote{},
		}

		normativeModelToReturn := &models.Normative{
			Key:           models.NormativeKey{ID: "norm-1", Version: 1},
			LawBucketID:   999,
			Type:          "article",
			Label:         "Article 5",
			EffectiveDate: effectiveDate,
		}

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			getVotingRoundByIDFunc: func(ctx context.Context, votingRoundID int) (*models.VotingRound, error) {
				return votingRoundModelToReturn, nil
			},
		}

		fakeNormativeRepositoryInstance := &fakeNormativeRepository{
			getNormativeByIDAndVersionFunc: func(ctx context.Context, normativeID string, normativeVersion int) (*models.Normative, error) {
				return normativeModelToReturn, nil
			},
		}

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{
			getLawBucketByIDFunc: func(ctx context.Context, lawBucketID int) (*models.LawBucket, error) {
				return nil, repository.ErrNotFound
			},
		}

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		votingServiceInstance := service.NewVotingService(
			fakeVotingRoundRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeLawBucketRepositoryInstance,
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
		)

		votingRoundDetailDTOResult, errResult := votingServiceInstance.GetVotingRoundByID(testContext, 1)

		require.NoError(testRunner, errResult)
		require.NotNil(testRunner, votingRoundDetailDTOResult)
		require.Equal(testRunner, 1, votingRoundDetailDTOResult.Id)
		// Normative should be populated
		require.Equal(testRunner, "norm-1", votingRoundDetailDTOResult.Normative.Id)
		// Law bucket should be zero-valued because lookup failed
		require.Equal(testRunner, 0, votingRoundDetailDTOResult.LawBucket.Id)
	})

	testRunner.Run("propagates ErrNotFound when round itself is missing", func(testRunner *testing.T) {
		testContext := context.Background()

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			getVotingRoundByIDFunc: func(ctx context.Context, votingRoundID int) (*models.VotingRound, error) {
				return nil, repository.ErrNotFound
			},
		}

		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{}
		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		votingServiceInstance := service.NewVotingService(
			fakeVotingRoundRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeLawBucketRepositoryInstance,
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
		)

		_, errResult := votingServiceInstance.GetVotingRoundByID(testContext, 999)

		require.Error(testRunner, errResult)
		require.True(testRunner, errors.Is(errResult, repository.ErrNotFound))
	})
}

func TestVotingServiceGetVotesByVotingRound(testRunner *testing.T) {
	testRunner.Run("batches politician and party lookups: calls each repository exactly once", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")

		votingRoundModelToReturn := &models.VotingRound{
			ID:               1,
			Title:            "Round 1",
			VoteDate:         voteDate,
			NormativeID:      "norm-1",
			NormativeVersion: 1,
			Votes: []models.Vote{
				{PoliticianID: "pol-1", PartyID: "party-1", Value: "Yes"},
				{PoliticianID: "pol-2", PartyID: "party-2", Value: "No"},
				{PoliticianID: "pol-3", PartyID: "party-1", Value: "Abstain"},
				{PoliticianID: "pol-1", PartyID: "party-1", Value: "Yes"},
				{PoliticianID: "pol-2", PartyID: "party-2", Value: "No"},
			},
		}

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticiansByIDsFunc: func(ctx context.Context, politicianIDs []string) (map[string]models.Politician, error) {
				result := make(map[string]models.Politician)
				for _, politicianID := range politicianIDs {
					result[politicianID] = models.Politician{ID: politicianID, Name: "Politician " + politicianID}
				}
				return result, nil
			},
		}

		fakePartyRepositoryInstance := &fakePartyRepository{
			getPartiesByIDsFunc: func(ctx context.Context, partyIDs []string) (map[string]models.Party, error) {
				result := make(map[string]models.Party)
				for _, partyID := range partyIDs {
					result[partyID] = models.Party{ID: partyID, Name: "Party " + partyID}
				}
				return result, nil
			},
		}

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			getVotingRoundByIDFunc: func(ctx context.Context, votingRoundID int) (*models.VotingRound, error) {
				return votingRoundModelToReturn, nil
			},
		}

		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{}

		votingServiceInstance := service.NewVotingService(
			fakeVotingRoundRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeLawBucketRepositoryInstance,
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
		)

		hydratedVotesResult, errResult := votingServiceInstance.GetVotesByVotingRound(testContext, 1)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, hydratedVotesResult, 5)
		// Verify batching: each repository called exactly once despite 5 votes across 2 parties / 3 politicians
		require.Equal(testRunner, 1, fakePoliticianRepositoryInstance.getPoliticiansByIDsCallCount,
			"politician repository should be called exactly once for batching")
		require.Equal(testRunner, 1, fakePartyRepositoryInstance.getPartiesByIDsCallCount,
			"party repository should be called exactly once for batching")
	})

	testRunner.Run("preserves order of embedded votes array", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")

		votingRoundModelToReturn := &models.VotingRound{
			ID:               1,
			VoteDate:         voteDate,
			NormativeID:      "norm-1",
			NormativeVersion: 1,
			Votes: []models.Vote{
				{PoliticianID: "pol-3", PartyID: "party-1", Value: "Yes"},
				{PoliticianID: "pol-1", PartyID: "party-1", Value: "No"},
				{PoliticianID: "pol-2", PartyID: "party-2", Value: "Abstain"},
			},
		}

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticiansByIDsFunc: func(ctx context.Context, politicianIDs []string) (map[string]models.Politician, error) {
				result := make(map[string]models.Politician)
				for _, politicianID := range politicianIDs {
					result[politicianID] = models.Politician{ID: politicianID}
				}
				return result, nil
			},
		}

		fakePartyRepositoryInstance := &fakePartyRepository{
			getPartiesByIDsFunc: func(ctx context.Context, partyIDs []string) (map[string]models.Party, error) {
				result := make(map[string]models.Party)
				for _, partyID := range partyIDs {
					result[partyID] = models.Party{ID: partyID}
				}
				return result, nil
			},
		}

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			getVotingRoundByIDFunc: func(ctx context.Context, votingRoundID int) (*models.VotingRound, error) {
				return votingRoundModelToReturn, nil
			},
		}

		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{}

		votingServiceInstance := service.NewVotingService(
			fakeVotingRoundRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeLawBucketRepositoryInstance,
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
		)

		hydratedVotesResult, errResult := votingServiceInstance.GetVotesByVotingRound(testContext, 1)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, hydratedVotesResult, 3)
		require.Equal(testRunner, "pol-3", hydratedVotesResult[0].PoliticianId)
		require.Equal(testRunner, "pol-1", hydratedVotesResult[1].PoliticianId)
		require.Equal(testRunner, "pol-2", hydratedVotesResult[2].PoliticianId)
	})

	testRunner.Run("handles missing politician/party in maps without panicking", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")

		votingRoundModelToReturn := &models.VotingRound{
			ID:               1,
			VoteDate:         voteDate,
			NormativeID:      "norm-1",
			NormativeVersion: 1,
			Votes: []models.Vote{
				{PoliticianID: "pol-1", PartyID: "party-1", Value: "Yes"},
				{PoliticianID: "pol-unknown", PartyID: "party-2", Value: "No"},
			},
		}

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticiansByIDsFunc: func(ctx context.Context, politicianIDs []string) (map[string]models.Politician, error) {
				// Return only pol-1, not pol-unknown
				return map[string]models.Politician{
					"pol-1": {ID: "pol-1", Name: "Known Politician"},
				}, nil
			},
		}

		fakePartyRepositoryInstance := &fakePartyRepository{
			getPartiesByIDsFunc: func(ctx context.Context, partyIDs []string) (map[string]models.Party, error) {
				result := make(map[string]models.Party)
				for _, partyID := range partyIDs {
					result[partyID] = models.Party{ID: partyID, Name: "Party " + partyID}
				}
				return result, nil
			},
		}

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			getVotingRoundByIDFunc: func(ctx context.Context, votingRoundID int) (*models.VotingRound, error) {
				return votingRoundModelToReturn, nil
			},
		}

		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{}

		votingServiceInstance := service.NewVotingService(
			fakeVotingRoundRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeLawBucketRepositoryInstance,
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
		)

		hydratedVotesResult, errResult := votingServiceInstance.GetVotesByVotingRound(testContext, 1)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, hydratedVotesResult, 2)

		// First vote has politician
		require.Equal(testRunner, "pol-1", hydratedVotesResult[0].PoliticianId)
		require.Equal(testRunner, "Known Politician", hydratedVotesResult[0].Politician.Name)

		// Second vote's politician not in map, should be zero-valued
		require.Equal(testRunner, "pol-unknown", hydratedVotesResult[1].PoliticianId)
		require.Equal(testRunner, "", hydratedVotesResult[1].Politician.Id)
	})

	testRunner.Run("propagates ErrNotFound for missing round", func(testRunner *testing.T) {
		testContext := context.Background()

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			getVotingRoundByIDFunc: func(ctx context.Context, votingRoundID int) (*models.VotingRound, error) {
				return nil, repository.ErrNotFound
			},
		}

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakePartyRepositoryInstance := &fakePartyRepository{}
		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{}

		votingServiceInstance := service.NewVotingService(
			fakeVotingRoundRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeLawBucketRepositoryInstance,
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
		)

		_, errResult := votingServiceInstance.GetVotesByVotingRound(testContext, 999)

		require.Error(testRunner, errResult)
		require.True(testRunner, errors.Is(errResult, repository.ErrNotFound))
	})
}
