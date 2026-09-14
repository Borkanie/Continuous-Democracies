package service

import (
	"context"

	"github.com/borkanie/brand-new-day-api/internal/generated"
	"github.com/borkanie/brand-new-day-api/internal/repository"
)

// PoliticianService implements the business logic for politicians, parties,
// and the enriched "votes cast by a politician" aggregation.
type PoliticianService struct {
	politicianRepository  repository.PoliticianRepository
	partyRepository       repository.PartyRepository
	votingRoundRepository repository.VotingRoundRepository
}

// NewPoliticianService constructs a PoliticianService from its repository dependencies.
func NewPoliticianService(
	politicianRepository repository.PoliticianRepository,
	partyRepository repository.PartyRepository,
	votingRoundRepository repository.VotingRoundRepository,
) *PoliticianService {
	return &PoliticianService{
		politicianRepository:  politicianRepository,
		partyRepository:       partyRepository,
		votingRoundRepository: votingRoundRepository,
	}
}

// ListParties returns every party.
func (service *PoliticianService) ListParties(requestContext context.Context) ([]generated.Party, error) {
	partyModels, err := service.partyRepository.ListParties(requestContext)
	if err != nil {
		return nil, err
	}

	partyDTOs := make([]generated.Party, 0, len(partyModels))
	for _, partyModel := range partyModels {
		partyDTOs = append(partyDTOs, mapPartyModelToDTO(partyModel))
	}
	return partyDTOs, nil
}

// GetPartyByID returns a single party by id, or repository.ErrNotFound if absent.
func (service *PoliticianService) GetPartyByID(requestContext context.Context, partyID string) (*generated.Party, error) {
	partyModel, err := service.partyRepository.GetPartyByID(requestContext, partyID)
	if err != nil {
		return nil, err
	}

	partyDTO := mapPartyModelToDTO(*partyModel)
	return &partyDTO, nil
}

// ListPoliticians returns every politician.
func (service *PoliticianService) ListPoliticians(requestContext context.Context) ([]generated.Politician, error) {
	politicianModels, err := service.politicianRepository.ListPoliticians(requestContext)
	if err != nil {
		return nil, err
	}

	politicianDTOs := make([]generated.Politician, 0, len(politicianModels))
	for _, politicianModel := range politicianModels {
		politicianDTOs = append(politicianDTOs, mapPoliticianModelToDTO(politicianModel))
	}
	return politicianDTOs, nil
}

// GetPoliticianByID returns a single politician by id, or repository.ErrNotFound if absent.
func (service *PoliticianService) GetPoliticianByID(requestContext context.Context, politicianID string) (*generated.Politician, error) {
	politicianModel, err := service.politicianRepository.GetPoliticianByID(requestContext, politicianID)
	if err != nil {
		return nil, err
	}

	politicianDTO := mapPoliticianModelToDTO(*politicianModel)
	return &politicianDTO, nil
}

// GetVotesByPolitician returns every vote cast by the given politician, each enriched with the
// resolved voting round, normative, and law bucket. It first confirms the politician exists so
// an unknown id yields repository.ErrNotFound (-> HTTP 404) rather than an empty 200.
func (service *PoliticianService) GetVotesByPolitician(requestContext context.Context, politicianID string) ([]generated.PoliticianVoteEntry, error) {
	if _, err := service.politicianRepository.GetPoliticianByID(requestContext, politicianID); err != nil {
		return nil, err
	}

	voteRecords, err := service.votingRoundRepository.ListVoteRecordsByPoliticianID(requestContext, politicianID)
	if err != nil {
		return nil, err
	}

	voteEntryDTOs := make([]generated.PoliticianVoteEntry, 0, len(voteRecords))
	for _, voteRecord := range voteRecords {
		voteEntryDTO := generated.PoliticianVoteEntry{
			Value:       generated.VoteValue(voteRecord.Value),
			VotingRound: mapVotingRoundModelToDTO(voteRecord.VotingRound),
		}
		if voteRecord.Normative != nil {
			voteEntryDTO.Normative = mapNormativeModelToDTO(*voteRecord.Normative)
		}
		if voteRecord.LawBucket != nil {
			voteEntryDTO.LawBucket = mapLawBucketModelToDTO(*voteRecord.LawBucket)
		}
		voteEntryDTOs = append(voteEntryDTOs, voteEntryDTO)
	}
	return voteEntryDTOs, nil
}
