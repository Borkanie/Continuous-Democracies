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

func TestPoliticianServiceListParties(testRunner *testing.T) {
	testRunner.Run("returns parties mapped correctly from models to DTOs", func(testRunner *testing.T) {
		testContext := context.Background()
		expectedPartyModels := []models.Party{
			{
				ID:      "party-1",
				Name:    "Party One",
				Acronym: "P1",
				LogoURL: "https://example.com/logo1.png",
				Color:   "#FF0000",
				Active:  true,
			},
			{
				ID:      "party-2",
				Name:    "Party Two",
				Acronym: "P2",
				LogoURL: "https://example.com/logo2.png",
				Color:   "#00FF00",
				Active:  false,
			},
		}

		fakePartyRepositoryInstance := &fakePartyRepository{
			listPartiesResult: expectedPartyModels,
		}
		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		partyDTOsResult, errResult := politicianServiceInstance.ListParties(testContext)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, partyDTOsResult, 2)

		require.Equal(testRunner, "party-1", partyDTOsResult[0].Id)
		require.Equal(testRunner, "Party One", partyDTOsResult[0].Name)
		require.Equal(testRunner, "P1", partyDTOsResult[0].Acronym)
		require.Equal(testRunner, "https://example.com/logo1.png", partyDTOsResult[0].LogoUrl)
		require.Equal(testRunner, "#FF0000", partyDTOsResult[0].Color)
		require.Equal(testRunner, true, partyDTOsResult[0].Active)

		require.Equal(testRunner, "party-2", partyDTOsResult[1].Id)
		require.Equal(testRunner, "Party Two", partyDTOsResult[1].Name)
		require.Equal(testRunner, "P2", partyDTOsResult[1].Acronym)
	})

	testRunner.Run("propagates repository error", func(testRunner *testing.T) {
		testContext := context.Background()
		expectedError := errors.New("database connection failed")

		fakePartyRepositoryInstance := &fakePartyRepository{
			listPartiesError: expectedError,
		}
		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		_, errResult := politicianServiceInstance.ListParties(testContext)

		require.ErrorIs(testRunner, errResult, expectedError)
	})
}

func TestPoliticianServiceGetPartyByID(testRunner *testing.T) {
	testRunner.Run("returns party mapped correctly from model to DTO", func(testRunner *testing.T) {
		testContext := context.Background()
		expectedPartyModel := &models.Party{
			ID:      "party-abc",
			Name:    "Test Party",
			Acronym: "TP",
			LogoURL: "https://example.com/test-logo.png",
			Color:   "#0000FF",
			Active:  true,
		}

		fakePartyRepositoryInstance := &fakePartyRepository{
			getPartyByIDFunc: func(ctx context.Context, partyID string) (*models.Party, error) {
				if partyID == "party-abc" {
					return expectedPartyModel, nil
				}
				return nil, repository.ErrNotFound
			},
		}
		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		partyDTOResult, errResult := politicianServiceInstance.GetPartyByID(testContext, "party-abc")

		require.NoError(testRunner, errResult)
		require.NotNil(testRunner, partyDTOResult)
		require.Equal(testRunner, "party-abc", partyDTOResult.Id)
		require.Equal(testRunner, "Test Party", partyDTOResult.Name)
		require.Equal(testRunner, "TP", partyDTOResult.Acronym)
		require.Equal(testRunner, "https://example.com/test-logo.png", partyDTOResult.LogoUrl)
	})

	testRunner.Run("propagates ErrNotFound from repository", func(testRunner *testing.T) {
		testContext := context.Background()

		fakePartyRepositoryInstance := &fakePartyRepository{
			getPartyByIDFunc: func(ctx context.Context, partyID string) (*models.Party, error) {
				return nil, repository.ErrNotFound
			},
		}
		fakePoliticianRepositoryInstance := &fakePoliticianRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		_, errResult := politicianServiceInstance.GetPartyByID(testContext, "nonexistent-party")

		require.Error(testRunner, errResult)
		require.True(testRunner, errors.Is(errResult, repository.ErrNotFound), "error should be ErrNotFound")
	})
}

func TestPoliticianServiceListPoliticians(testRunner *testing.T) {
	testRunner.Run("returns politicians mapped correctly from models to DTOs", func(testRunner *testing.T) {
		testContext := context.Background()
		expectedPoliticianModels := []models.Politician{
			{
				ID:           "pol-1",
				Name:         "John Doe",
				Gender:       1,
				ImageURL:     "https://example.com/john.png",
				PartyID:      "party-1",
				Active:       true,
				WorkLocation: 1,
			},
			{
				ID:           "pol-2",
				Name:         "Jane Smith",
				Gender:       2,
				ImageURL:     "https://example.com/jane.png",
				PartyID:      "party-2",
				Active:       false,
				WorkLocation: 0,
			},
		}

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			listPoliticiansResult: expectedPoliticianModels,
		}
		fakePartyRepositoryInstance := &fakePartyRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		politicianDTOsResult, errResult := politicianServiceInstance.ListPoliticians(testContext)

		require.NoError(testRunner, errResult)
		require.Len(testRunner, politicianDTOsResult, 2)

		require.Equal(testRunner, "pol-1", politicianDTOsResult[0].Id)
		require.Equal(testRunner, "John Doe", politicianDTOsResult[0].Name)
		require.Equal(testRunner, 1, politicianDTOsResult[0].Gender)
		require.Equal(testRunner, "https://example.com/john.png", politicianDTOsResult[0].ImageUrl)
		require.Equal(testRunner, "party-1", politicianDTOsResult[0].PartyId)
		require.Equal(testRunner, true, politicianDTOsResult[0].Active)
		require.Equal(testRunner, 1, politicianDTOsResult[0].WorkLocation)

		require.Equal(testRunner, "pol-2", politicianDTOsResult[1].Id)
		require.Equal(testRunner, "Jane Smith", politicianDTOsResult[1].Name)
	})
}

func TestPoliticianServiceGetPoliticianByID(testRunner *testing.T) {
	testRunner.Run("returns politician mapped correctly from model to DTO", func(testRunner *testing.T) {
		testContext := context.Background()
		expectedPoliticianModel := &models.Politician{
			ID:           "pol-xyz",
			Name:         "Test Politician",
			Gender:       1,
			ImageURL:     "https://example.com/test.png",
			PartyID:      "party-xyz",
			Active:       true,
			WorkLocation: 1,
		}

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticianByIDFunc: func(ctx context.Context, politicianID string) (*models.Politician, error) {
				if politicianID == "pol-xyz" {
					return expectedPoliticianModel, nil
				}
				return nil, repository.ErrNotFound
			},
		}
		fakePartyRepositoryInstance := &fakePartyRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		politicianDTOResult, errResult := politicianServiceInstance.GetPoliticianByID(testContext, "pol-xyz")

		require.NoError(testRunner, errResult)
		require.NotNil(testRunner, politicianDTOResult)
		require.Equal(testRunner, "pol-xyz", politicianDTOResult.Id)
		require.Equal(testRunner, "Test Politician", politicianDTOResult.Name)
		require.Equal(testRunner, "https://example.com/test.png", politicianDTOResult.ImageUrl)
	})

	testRunner.Run("propagates ErrNotFound from repository", func(testRunner *testing.T) {
		testContext := context.Background()

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticianByIDFunc: func(ctx context.Context, politicianID string) (*models.Politician, error) {
				return nil, repository.ErrNotFound
			},
		}
		fakePartyRepositoryInstance := &fakePartyRepository{}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		_, errResult := politicianServiceInstance.GetPoliticianByID(testContext, "nonexistent")

		require.Error(testRunner, errResult)
		require.True(testRunner, errors.Is(errResult, repository.ErrNotFound))
	})
}

func TestPoliticianServiceGetVotesByPolitician(testRunner *testing.T) {
	testRunner.Run("returns empty list when politician exists and has no votes", func(testRunner *testing.T) {
		testContext := context.Background()

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticianByIDFunc: func(ctx context.Context, politicianID string) (*models.Politician, error) {
				return &models.Politician{ID: "pol-1"}, nil
			},
		}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			listVoteRecordsByPoliticianIDFunc: func(ctx context.Context, politicianID string) ([]models.PoliticianVoteRecord, error) {
				return []models.PoliticianVoteRecord{}, nil
			},
		}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		votesResult, errResult := politicianServiceInstance.GetVotesByPolitician(testContext, "pol-1")

		require.NoError(testRunner, errResult)
		require.NotNil(testRunner, votesResult)
		require.Len(testRunner, votesResult, 0)
	})

	testRunner.Run("returns ErrNotFound when politician does not exist without calling voting repository", func(testRunner *testing.T) {
		testContext := context.Background()

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticianByIDFunc: func(ctx context.Context, politicianID string) (*models.Politician, error) {
				return nil, repository.ErrNotFound
			},
		}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		_, errResult := politicianServiceInstance.GetVotesByPolitician(testContext, "nonexistent")

		require.Error(testRunner, errResult)
		require.True(testRunner, errors.Is(errResult, repository.ErrNotFound))
		require.Equal(testRunner, 0, fakeVotingRoundRepositoryInstance.listVoteRecordsByPoliticianIDCallCount,
			"voting round repository should not be called when politician does not exist")
	})

	testRunner.Run("handles both normative and law bucket nil", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticianByIDFunc: func(ctx context.Context, politicianID string) (*models.Politician, error) {
				return &models.Politician{ID: "pol-1"}, nil
			},
		}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			listVoteRecordsByPoliticianIDFunc: func(ctx context.Context, politicianID string) ([]models.PoliticianVoteRecord, error) {
				return []models.PoliticianVoteRecord{
					{
						Value: "Yes",
						VotingRound: models.VotingRound{
							ID:               1,
							Title:            "Vote 1",
							NormativeID:      "norm-1",
							NormativeVersion: 1,
							VoteDate:         voteDate,
						},
						Normative: nil,
						LawBucket: nil,
					},
				}, nil
			},
		}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		votesResult, errResult := politicianServiceInstance.GetVotesByPolitician(testContext, "pol-1")

		require.NoError(testRunner, errResult)
		require.Len(testRunner, votesResult, 1)
		require.Equal(testRunner, "Yes", string(votesResult[0].Value))
		// Verify zero-valued DTOs when nil (not panicking)
		require.Equal(testRunner, "", votesResult[0].Normative.Id)
		require.Equal(testRunner, 0, votesResult[0].LawBucket.Id)
	})

	testRunner.Run("handles normative nil only", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticianByIDFunc: func(ctx context.Context, politicianID string) (*models.Politician, error) {
				return &models.Politician{ID: "pol-1"}, nil
			},
		}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			listVoteRecordsByPoliticianIDFunc: func(ctx context.Context, politicianID string) ([]models.PoliticianVoteRecord, error) {
				return []models.PoliticianVoteRecord{
					{
						Value: "No",
						VotingRound: models.VotingRound{
							ID:               1,
							NormativeID:      "norm-1",
							NormativeVersion: 1,
							VoteDate:         voteDate,
						},
						Normative: nil,
						LawBucket: &models.LawBucket{
							ID:       42,
							PLNumber: "PL 123/2024",
						},
					},
				}, nil
			},
		}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		votesResult, errResult := politicianServiceInstance.GetVotesByPolitician(testContext, "pol-1")

		require.NoError(testRunner, errResult)
		require.Len(testRunner, votesResult, 1)
		require.Equal(testRunner, "No", string(votesResult[0].Value))
		// Normative should be zero-valued
		require.Equal(testRunner, "", votesResult[0].Normative.Id)
		// Law bucket should be populated
		require.Equal(testRunner, 42, votesResult[0].LawBucket.Id)
		require.Equal(testRunner, "PL 123/2024", votesResult[0].LawBucket.PlNumber)
	})

	testRunner.Run("handles law bucket nil only", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticianByIDFunc: func(ctx context.Context, politicianID string) (*models.Politician, error) {
				return &models.Politician{ID: "pol-1"}, nil
			},
		}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			listVoteRecordsByPoliticianIDFunc: func(ctx context.Context, politicianID string) ([]models.PoliticianVoteRecord, error) {
				return []models.PoliticianVoteRecord{
					{
						Value: "Abstain",
						VotingRound: models.VotingRound{
							ID:               1,
							NormativeID:      "norm-1",
							NormativeVersion: 1,
							VoteDate:         voteDate,
						},
						Normative: &models.Normative{
							Key:   models.NormativeKey{ID: "norm-1", Version: 1},
							Label: "Article 5",
							Type:  "article",
						},
						LawBucket: nil,
					},
				}, nil
			},
		}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		votesResult, errResult := politicianServiceInstance.GetVotesByPolitician(testContext, "pol-1")

		require.NoError(testRunner, errResult)
		require.Len(testRunner, votesResult, 1)
		require.Equal(testRunner, "Abstain", string(votesResult[0].Value))
		// Normative should be populated
		require.Equal(testRunner, "norm-1", votesResult[0].Normative.Id)
		require.Equal(testRunner, "Article 5", votesResult[0].Normative.Label)
		// Law bucket should be zero-valued
		require.Equal(testRunner, 0, votesResult[0].LawBucket.Id)
	})

	testRunner.Run("handles both normative and law bucket populated", func(testRunner *testing.T) {
		testContext := context.Background()
		voteDate, _ := time.Parse("2006-01-02", "2024-01-15")

		fakePoliticianRepositoryInstance := &fakePoliticianRepository{
			getPoliticianByIDFunc: func(ctx context.Context, politicianID string) (*models.Politician, error) {
				return &models.Politician{ID: "pol-1"}, nil
			},
		}
		fakeVotingRoundRepositoryInstance := &fakeVotingRoundRepository{
			listVoteRecordsByPoliticianIDFunc: func(ctx context.Context, politicianID string) ([]models.PoliticianVoteRecord, error) {
				return []models.PoliticianVoteRecord{
					{
						Value: "Absent",
						VotingRound: models.VotingRound{
							ID:               1,
							NormativeID:      "norm-1",
							NormativeVersion: 2,
							VoteDate:         voteDate,
						},
						Normative: &models.Normative{
							Key:   models.NormativeKey{ID: "norm-1", Version: 2},
							Label: "Amendment 3",
							Type:  "amendment",
						},
						LawBucket: &models.LawBucket{
							ID:       99,
							PLNumber: "PL 999/2024",
						},
					},
				}, nil
			},
		}
		fakePartyRepositoryInstance := &fakePartyRepository{}

		politicianServiceInstance := service.NewPoliticianService(
			fakePoliticianRepositoryInstance,
			fakePartyRepositoryInstance,
			fakeVotingRoundRepositoryInstance,
		)

		votesResult, errResult := politicianServiceInstance.GetVotesByPolitician(testContext, "pol-1")

		require.NoError(testRunner, errResult)
		require.Len(testRunner, votesResult, 1)
		require.Equal(testRunner, "Absent", string(votesResult[0].Value))
		require.Equal(testRunner, "norm-1", votesResult[0].Normative.Id)
		require.Equal(testRunner, 2, votesResult[0].Normative.Version)
		require.Equal(testRunner, "Amendment 3", votesResult[0].Normative.Label)
		require.Equal(testRunner, 99, votesResult[0].LawBucket.Id)
	})
}
