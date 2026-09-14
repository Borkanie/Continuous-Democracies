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

func TestLawServiceListLawBuckets(testRunner *testing.T) {
	testRunner.Run("maps law buckets with embedded normatives, does not call normative repository", func(testRunner *testing.T) {
		testContext := context.Background()
		initiationDate, _ := time.Parse("2006-01-02", "2024-01-01")
		effectiveDate, _ := time.Parse("2006-01-02", "2024-01-10")

		lawBucketModelsToReturn := []models.LawBucket{
			{
				ID:             1,
				Version:        1,
				PLNumber:       "PL 100/2024",
				Title:          "Law 1",
				Description:    "Description 1",
				InitiationDate: initiationDate,
				Status:         "pending",
				Normatives: []models.Normative{
					{
						Key:           models.NormativeKey{ID: "norm-1", Version: 1},
						LawBucketID:   1,
						Type:          "article",
						Label:         "Article 1",
						Text:          "Article text",
						EffectiveDate: effectiveDate,
					},
				},
			},
		}

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{
			listLawBucketsResult: lawBucketModelsToReturn,
		}
		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		lawServiceInstance := service.NewLawService(
			fakeLawBucketRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		lawBucketDTOsResult, errResult := lawServiceInstance.ListLawBuckets(testContext)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, lawBucketDTOsResult, 1)
		require.Equal(testRunner, 1, lawBucketDTOsResult[0].Id)
		require.Equal(testRunner, "PL 100/2024", lawBucketDTOsResult[0].PlNumber)
		require.Len(testRunner, lawBucketDTOsResult[0].Normatives, 1)
		require.Equal(testRunner, "norm-1", lawBucketDTOsResult[0].Normatives[0].Id)
		// Verify normative repository was not called
		require.Equal(testRunner, 0, fakeNormativeRepositoryInstance.listCurrentNormativesByLawBucketIDCallCount,
			"normative repository should not be called when normatives are already embedded")
	})

	testRunner.Run("returns empty slice when no law buckets exist", func(testRunner *testing.T) {
		testContext := context.Background()

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{
			listLawBucketsResult: []models.LawBucket{},
		}
		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		lawServiceInstance := service.NewLawService(
			fakeLawBucketRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		lawBucketDTOsResult, errResult := lawServiceInstance.ListLawBuckets(testContext)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, lawBucketDTOsResult, 0)
	})

	testRunner.Run("returns non-nil normatives slice even when bucket has no normatives", func(testRunner *testing.T) {
		testContext := context.Background()
		initiationDate, _ := time.Parse("2006-01-02", "2024-01-01")

		lawBucketModelsToReturn := []models.LawBucket{
			{
				ID:             42,
				Version:        1,
				PLNumber:       "PL 999/2024",
				Title:          "Law with no normatives",
				InitiationDate: initiationDate,
				Status:         "pending",
				Normatives:     []models.Normative{},
			},
		}

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{
			listLawBucketsResult: lawBucketModelsToReturn,
		}
		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		lawServiceInstance := service.NewLawService(
			fakeLawBucketRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		lawBucketDTOsResult, errResult := lawServiceInstance.ListLawBuckets(testContext)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, lawBucketDTOsResult, 1)
		require.NotNil(testRunner, lawBucketDTOsResult[0].Normatives)
		require.Len(testRunner, lawBucketDTOsResult[0].Normatives, 0)
	})
}

func TestLawServiceGetLawBucketByID(testRunner *testing.T) {
	testRunner.Run("returns law bucket with embedded normatives", func(testRunner *testing.T) {
		testContext := context.Background()
		initiationDate, _ := time.Parse("2006-01-02", "2024-01-01")
		effectiveDate, _ := time.Parse("2006-01-02", "2024-01-10")

		lawBucketModelToReturn := &models.LawBucket{
			ID:             42,
			Version:        2,
			PLNumber:       "PL 123/2024",
			Title:          "Test Law",
			Description:    "Test Description",
			InitiationDate: initiationDate,
			Status:         "approved",
			Normatives: []models.Normative{
				{
					Key:           models.NormativeKey{ID: "norm-1", Version: 2},
					LawBucketID:   42,
					Type:          "article",
					Label:         "Article 3",
					Text:          "Text of article 3",
					EffectiveDate: effectiveDate,
				},
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
		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		lawServiceInstance := service.NewLawService(
			fakeLawBucketRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		lawBucketDTOResult, errResult := lawServiceInstance.GetLawBucketByID(testContext, 42)

		require.NoError(testRunner, errResult)
		require.NotNil(testRunner, lawBucketDTOResult)
		require.Equal(testRunner, 42, lawBucketDTOResult.Id)
		require.Equal(testRunner, 2, lawBucketDTOResult.Version)
		require.Equal(testRunner, "PL 123/2024", lawBucketDTOResult.PlNumber)
		require.Len(testRunner, lawBucketDTOResult.Normatives, 1)
		require.Equal(testRunner, "norm-1", lawBucketDTOResult.Normatives[0].Id)
		require.Equal(testRunner, 2, lawBucketDTOResult.Normatives[0].Version)
	})

	testRunner.Run("propagates ErrNotFound from repository", func(testRunner *testing.T) {
		testContext := context.Background()

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{
			getLawBucketByIDFunc: func(ctx context.Context, lawBucketID int) (*models.LawBucket, error) {
				return nil, repository.ErrNotFound
			},
		}
		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		lawServiceInstance := service.NewLawService(
			fakeLawBucketRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		_, errResult := lawServiceInstance.GetLawBucketByID(testContext, 999)

		require.Error(testRunner, errResult)
		require.True(testRunner, errors.Is(errResult, repository.ErrNotFound))
	})
}

func TestLawServiceGetVotingRoundsByLawBucket(testRunner *testing.T) {
	testRunner.Run("returns voting rounds with empty votes slice when law bucket exists", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")

		lawBucketModelToReturn := &models.LawBucket{
			ID:         42,
			PLNumber:   "PL 123/2024",
			Normatives: []models.Normative{},
		}

		votingRoundModelsToReturn := []models.VotingRound{
			{
				ID:               1,
				Title:            "Round 1",
				VoteDate:         voteDate,
				NormativeID:      "norm-1",
				NormativeVersion: 1,
				Votes: []models.Vote{
					{PoliticianID: "pol-1", PartyID: "party-1", Value: "Yes"},
					{PoliticianID: "pol-2", PartyID: "party-2", Value: "No"},
				},
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

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			listVotingRoundsByLawBucketIDFunc: func(ctx context.Context, lawBucketID int) ([]models.VotingRound, error) {
				if lawBucketID == 42 {
					return votingRoundModelsToReturn, nil
				}
				return []models.VotingRound{}, nil
			},
		}

		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}

		lawServiceInstance := service.NewLawService(
			fakeLawBucketRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		votingRoundDTOsResult, errResult := lawServiceInstance.GetVotingRoundsByLawBucket(testContext, 42)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, votingRoundDTOsResult, 1)
		require.Equal(testRunner, 1, votingRoundDTOsResult[0].Id)
		require.NotNil(testRunner, votingRoundDTOsResult[0].Votes)
		require.Len(testRunner, votingRoundDTOsResult[0].Votes, 0)
	})

	testRunner.Run("returns ErrNotFound when law bucket does not exist without calling voting round repository", func(testRunner *testing.T) {
		testContext := context.Background()

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{
			getLawBucketByIDFunc: func(ctx context.Context, lawBucketID int) (*models.LawBucket, error) {
				return nil, repository.ErrNotFound
			},
		}

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}
		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}

		lawServiceInstance := service.NewLawService(
			fakeLawBucketRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		_, errResult := lawServiceInstance.GetVotingRoundsByLawBucket(testContext, 999)

		require.Error(testRunner, errResult)
		require.True(testRunner, errors.Is(errResult, repository.ErrNotFound))
		require.Equal(testRunner, 0, fakeVotingRoundRepositoryInstance.listVotingRoundsByLawBucketIDCallCount,
			"voting round repository should not be called when law bucket does not exist")
	})

	testRunner.Run("returns empty slice when law bucket has no voting rounds", func(testRunner *testing.T) {
		testContext := context.Background()

		lawBucketModelToReturn := &models.LawBucket{
			ID:         42,
			PLNumber:   "PL 999/2024",
			Normatives: []models.Normative{},
		}

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{
			getLawBucketByIDFunc: func(ctx context.Context, lawBucketID int) (*models.LawBucket, error) {
				if lawBucketID == 42 {
					return lawBucketModelToReturn, nil
				}
				return nil, repository.ErrNotFound
			},
		}

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			listVotingRoundsByLawBucketIDFunc: func(ctx context.Context, lawBucketID int) ([]models.VotingRound, error) {
				return []models.VotingRound{}, nil
			},
		}

		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}

		lawServiceInstance := service.NewLawService(
			fakeLawBucketRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		votingRoundDTOsResult, errResult := lawServiceInstance.GetVotingRoundsByLawBucket(testContext, 42)

		require.NoError(testRunner, errResult)
		require.NotNil(testRunner, votingRoundDTOsResult)
		require.Len(testRunner, votingRoundDTOsResult, 0)
	})

	testRunner.Run("does not call normative repository", func(testRunner *testing.T) {
		testContext := context.Background()

		lawBucketModelToReturn := &models.LawBucket{
			ID:         42,
			Normatives: []models.Normative{},
		}

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{
			getLawBucketByIDFunc: func(ctx context.Context, lawBucketID int) (*models.LawBucket, error) {
				return lawBucketModelToReturn, nil
			},
		}

		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			listVotingRoundsByLawBucketIDFunc: func(ctx context.Context, lawBucketID int) ([]models.VotingRound, error) {
				return []models.VotingRound{}, nil
			},
		}

		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}

		lawServiceInstance := service.NewLawService(
			fakeLawBucketRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		_, errResult := lawServiceInstance.GetVotingRoundsByLawBucket(testContext, 42)

		require.NoError(testRunner, errResult)
		require.Equal(testRunner, 0, fakeNormativeRepositoryInstance.getNormativeByIDAndVersionCallCount,
			"normative repository should not be called")
		require.Equal(testRunner, 0, fakeNormativeRepositoryInstance.listCurrentNormativesByLawBucketIDCallCount,
			"normative repository should not be called")
	})
}

func TestLawServiceNormativeMapping(testRunner *testing.T) {
	testRunner.Run("maps normative key (ID, Version) to flat DTO fields correctly", func(testRunner *testing.T) {
		testContext := context.Background()
		initiationDate, _ := time.Parse("2006-01-02", "2024-01-01")
		effectiveDate, _ := time.Parse("2006-01-02", "2024-01-10")

		lawBucketModelToReturn := &models.LawBucket{
			ID:             42,
			Version:        1,
			PLNumber:       "PL 123/2024",
			Title:          "Test Law",
			InitiationDate: initiationDate,
			Status:         "pending",
			Normatives: []models.Normative{
				{
					Key:           models.NormativeKey{ID: "norm-abc-123", Version: 5},
					LawBucketID:   42,
					Type:          "article",
					Label:         "Article 10",
					Text:          "Article text",
					EffectiveDate: effectiveDate,
				},
			},
		}

		fakeLawBucketRepositoryInstance := &fakeLawBucketRepository{
			getLawBucketByIDFunc: func(ctx context.Context, lawBucketID int) (*models.LawBucket, error) {
				return lawBucketModelToReturn, nil
			},
		}
		fakeNormativeRepositoryInstance := &fakeNormativeRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		lawServiceInstance := service.NewLawService(
			fakeLawBucketRepositoryInstance,
			fakeNormativeRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		lawBucketDTOResult, errResult := lawServiceInstance.GetLawBucketByID(testContext, 42)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, lawBucketDTOResult.Normatives, 1)

		normativeDTO := lawBucketDTOResult.Normatives[0]
		require.Equal(testRunner, "norm-abc-123", normativeDTO.Id, "ID should be mapped from Key.ID")
		require.Equal(testRunner, 5, normativeDTO.Version, "Version should be mapped from Key.Version")
		require.Equal(testRunner, "Article 10", normativeDTO.Label)
		require.Equal(testRunner, "article", string(normativeDTO.Type))
	})
}
