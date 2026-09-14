package service

import (
	"context"

	"github.com/borkanie/brand-new-day-api/internal/generated"
	"github.com/borkanie/brand-new-day-api/internal/repository"
)

// LawService implements the business logic for law buckets, including resolving the voting
// rounds associated with a law bucket. The repository layer already embeds each law bucket's
// current normatives, so this service does not need to re-fetch them.
type LawService struct {
	lawBucketRepository   repository.LawBucketRepository
	normativeRepository   repository.NormativeRepository
	votingRoundRepository repository.VotingRoundRepository
}

// NewLawService constructs a LawService from its repository dependencies.
func NewLawService(
	lawBucketRepository repository.LawBucketRepository,
	normativeRepository repository.NormativeRepository,
	votingRoundRepository repository.VotingRoundRepository,
) *LawService {
	return &LawService{
		lawBucketRepository:   lawBucketRepository,
		normativeRepository:   normativeRepository,
		votingRoundRepository: votingRoundRepository,
	}
}

// ListLawBuckets returns every law bucket, with its current normatives embedded.
func (service *LawService) ListLawBuckets(requestContext context.Context) ([]generated.LawBucket, error) {
	lawBucketModels, err := service.lawBucketRepository.ListLawBuckets(requestContext)
	if err != nil {
		return nil, err
	}

	lawBucketDTOs := make([]generated.LawBucket, 0, len(lawBucketModels))
	for _, lawBucketModel := range lawBucketModels {
		lawBucketDTOs = append(lawBucketDTOs, mapLawBucketModelToDTO(lawBucketModel))
	}
	return lawBucketDTOs, nil
}

// GetLawBucketByID returns a single law bucket by id, with its current normatives embedded,
// or repository.ErrNotFound if absent.
func (service *LawService) GetLawBucketByID(requestContext context.Context, lawBucketID int) (*generated.LawBucket, error) {
	lawBucketModel, err := service.lawBucketRepository.GetLawBucketByID(requestContext, lawBucketID)
	if err != nil {
		return nil, err
	}

	lawBucketDTO := mapLawBucketModelToDTO(*lawBucketModel)
	return &lawBucketDTO, nil
}

// GetVotingRoundsByLawBucket returns every voting round associated with the given law bucket.
// It first confirms the law bucket exists so an unknown id yields repository.ErrNotFound
// (-> HTTP 404) rather than an empty 200. Each round's Votes field is emptied to keep the list
// response small, same rationale as VotingService.ListVotingRounds.
func (service *LawService) GetVotingRoundsByLawBucket(requestContext context.Context, lawBucketID int) ([]generated.VotingRound, error) {
	if _, err := service.lawBucketRepository.GetLawBucketByID(requestContext, lawBucketID); err != nil {
		return nil, err
	}

	votingRoundModels, err := service.votingRoundRepository.ListVotingRoundsByLawBucketID(requestContext, lawBucketID)
	if err != nil {
		return nil, err
	}

	return mapVotingRoundsModelToDTOWithEmptyVotes(votingRoundModels), nil
}
