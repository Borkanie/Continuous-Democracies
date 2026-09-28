package service_test

import (
	"context"

	"github.com/borkanie/brand-new-day-api/internal/models"
	"github.com/borkanie/brand-new-day-api/internal/repository"
)

// fakePoliticianRepository is a hand-written test double for repository.PoliticianRepository.
type fakePoliticianRepository struct {
	listPoliticiansResult        []models.Politician
	listPoliticiansError         error
	getPoliticianByIDFunc        func(ctx context.Context, politicianID string) (*models.Politician, error)
	getPoliticiansByIDsFunc      func(ctx context.Context, politicianIDs []string) (map[string]models.Politician, error)
	getPoliticianByIDCallCount   int
	getPoliticiansByIDsCallCount int
}

func (fakePolitician *fakePoliticianRepository) ListPoliticians(requestContext context.Context) ([]models.Politician, error) {
	return fakePolitician.listPoliticiansResult, fakePolitician.listPoliticiansError
}

func (fakePolitician *fakePoliticianRepository) GetPoliticianByID(requestContext context.Context, politicianID string) (*models.Politician, error) {
	fakePolitician.getPoliticianByIDCallCount++
	if fakePolitician.getPoliticianByIDFunc != nil {
		return fakePolitician.getPoliticianByIDFunc(requestContext, politicianID)
	}
	return nil, repository.ErrNotFound
}

func (fakePolitician *fakePoliticianRepository) GetPoliticiansByIDs(requestContext context.Context, politicianIDs []string) (map[string]models.Politician, error) {
	fakePolitician.getPoliticiansByIDsCallCount++
	if fakePolitician.getPoliticiansByIDsFunc != nil {
		return fakePolitician.getPoliticiansByIDsFunc(requestContext, politicianIDs)
	}
	return make(map[string]models.Politician), nil
}

// fakePartyRepository is a hand-written test double for repository.PartyRepository.
type fakePartyRepository struct {
	listPartiesResult        []models.Party
	listPartiesError         error
	getPartyByIDFunc         func(ctx context.Context, partyID string) (*models.Party, error)
	getPartiesByIDsFunc      func(ctx context.Context, partyIDs []string) (map[string]models.Party, error)
	getPartyByIDCallCount    int
	getPartiesByIDsCallCount int
}

func (fakeParty *fakePartyRepository) ListParties(requestContext context.Context) ([]models.Party, error) {
	return fakeParty.listPartiesResult, fakeParty.listPartiesError
}

func (fakeParty *fakePartyRepository) GetPartyByID(requestContext context.Context, partyID string) (*models.Party, error) {
	fakeParty.getPartyByIDCallCount++
	if fakeParty.getPartyByIDFunc != nil {
		return fakeParty.getPartyByIDFunc(requestContext, partyID)
	}
	return nil, repository.ErrNotFound
}

func (fakeParty *fakePartyRepository) GetPartiesByIDs(requestContext context.Context, partyIDs []string) (map[string]models.Party, error) {
	fakeParty.getPartiesByIDsCallCount++
	if fakeParty.getPartiesByIDsFunc != nil {
		return fakeParty.getPartiesByIDsFunc(requestContext, partyIDs)
	}
	return make(map[string]models.Party), nil
}

// fakeVotingRoundRepository is a hand-written test double for repository.VotingRoundRepository.
type fakeVotingRoundRepository struct {
	listVotingRoundsResult                 []models.VotingRound
	listVotingRoundsError                  error
	getVotingRoundByIDFunc                 func(ctx context.Context, votingRoundID int) (*models.VotingRound, error)
	listVotingRoundsByLawBucketIDFunc      func(ctx context.Context, lawBucketID int) ([]models.VotingRound, error)
	listVoteRecordsByPoliticianIDFunc      func(ctx context.Context, politicianID string) ([]models.PoliticianVoteRecord, error)
	getVotingRoundByIDCallCount            int
	listVotingRoundsByLawBucketIDCallCount int
	listVoteRecordsByPoliticianIDCallCount int
}

func (fakeVotingRound *fakeVotingRoundRepository) ListVotingRounds(requestContext context.Context) ([]models.VotingRound, error) {
	return fakeVotingRound.listVotingRoundsResult, fakeVotingRound.listVotingRoundsError
}

func (fakeVotingRound *fakeVotingRoundRepository) GetVotingRoundByID(requestContext context.Context, votingRoundID int) (*models.VotingRound, error) {
	fakeVotingRound.getVotingRoundByIDCallCount++
	if fakeVotingRound.getVotingRoundByIDFunc != nil {
		return fakeVotingRound.getVotingRoundByIDFunc(requestContext, votingRoundID)
	}
	return nil, repository.ErrNotFound
}

func (fakeVotingRound *fakeVotingRoundRepository) ListVotingRoundsByLawBucketID(requestContext context.Context, lawBucketID int) ([]models.VotingRound, error) {
	fakeVotingRound.listVotingRoundsByLawBucketIDCallCount++
	if fakeVotingRound.listVotingRoundsByLawBucketIDFunc != nil {
		return fakeVotingRound.listVotingRoundsByLawBucketIDFunc(requestContext, lawBucketID)
	}
	return []models.VotingRound{}, nil
}

func (fakeVotingRound *fakeVotingRoundRepository) ListVoteRecordsByPoliticianID(requestContext context.Context, politicianID string) ([]models.PoliticianVoteRecord, error) {
	fakeVotingRound.listVoteRecordsByPoliticianIDCallCount++
	if fakeVotingRound.listVoteRecordsByPoliticianIDFunc != nil {
		return fakeVotingRound.listVoteRecordsByPoliticianIDFunc(requestContext, politicianID)
	}
	return []models.PoliticianVoteRecord{}, nil
}

// fakeNormativeRepository is a hand-written test double for repository.NormativeRepository.
type fakeNormativeRepository struct {
	getNormativeByIDAndVersionFunc              func(ctx context.Context, normativeID string, normativeVersion int) (*models.Normative, error)
	listCurrentNormativesByLawBucketIDFunc      func(ctx context.Context, lawBucketID int) ([]models.Normative, error)
	getNormativeByIDAndVersionCallCount         int
	listCurrentNormativesByLawBucketIDCallCount int
}

func (fakeNormative *fakeNormativeRepository) GetNormativeByIDAndVersion(requestContext context.Context, normativeID string, normativeVersion int) (*models.Normative, error) {
	fakeNormative.getNormativeByIDAndVersionCallCount++
	if fakeNormative.getNormativeByIDAndVersionFunc != nil {
		return fakeNormative.getNormativeByIDAndVersionFunc(requestContext, normativeID, normativeVersion)
	}
	return nil, repository.ErrNotFound
}

func (fakeNormative *fakeNormativeRepository) ListCurrentNormativesByLawBucketID(requestContext context.Context, lawBucketID int) ([]models.Normative, error) {
	fakeNormative.listCurrentNormativesByLawBucketIDCallCount++
	if fakeNormative.listCurrentNormativesByLawBucketIDFunc != nil {
		return fakeNormative.listCurrentNormativesByLawBucketIDFunc(requestContext, lawBucketID)
	}
	return []models.Normative{}, nil
}

// fakeLawBucketRepository is a hand-written test double for repository.LawBucketRepository.
type fakeLawBucketRepository struct {
	listLawBucketsResult      []models.LawBucket
	listLawBucketsError       error
	getLawBucketByIDFunc      func(ctx context.Context, lawBucketID int) (*models.LawBucket, error)
	getLawBucketByIDCallCount int
}

func (fakeLawBucket *fakeLawBucketRepository) ListLawBuckets(requestContext context.Context) ([]models.LawBucket, error) {
	return fakeLawBucket.listLawBucketsResult, fakeLawBucket.listLawBucketsError
}

func (fakeLawBucket *fakeLawBucketRepository) GetLawBucketByID(requestContext context.Context, lawBucketID int) (*models.LawBucket, error) {
	fakeLawBucket.getLawBucketByIDCallCount++
	if fakeLawBucket.getLawBucketByIDFunc != nil {
		return fakeLawBucket.getLawBucketByIDFunc(requestContext, lawBucketID)
	}
	return nil, repository.ErrNotFound
}
