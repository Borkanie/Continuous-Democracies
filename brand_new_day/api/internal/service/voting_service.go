package service

import (
	"context"
	"log/slog"

	"github.com/borkanie/brand-new-day-api/internal/generated"
	"github.com/borkanie/brand-new-day-api/internal/repository"
)

// VotingService implements the business logic for voting rounds and the votes embedded on them,
// including hydration of those votes with the voting politician and party.
type VotingService struct {
	votingRoundRepository repository.VotingRoundRepository
	normativeRepository   repository.NormativeRepository
	lawBucketRepository   repository.LawBucketRepository
	politicianRepository  repository.PoliticianRepository
	partyRepository       repository.PartyRepository
}

// NewVotingService constructs a VotingService from its repository dependencies.
func NewVotingService(
	votingRoundRepository repository.VotingRoundRepository,
	normativeRepository repository.NormativeRepository,
	lawBucketRepository repository.LawBucketRepository,
	politicianRepository repository.PoliticianRepository,
	partyRepository repository.PartyRepository,
) *VotingService {
	return &VotingService{
		votingRoundRepository: votingRoundRepository,
		normativeRepository:   normativeRepository,
		lawBucketRepository:   lawBucketRepository,
		politicianRepository:  politicianRepository,
		partyRepository:       partyRepository,
	}
}

// ListVotingRounds returns every voting round. Each round's Votes field is emptied to keep the
// list response small; fetch GET /votingRounds/{roundId}/votes for the hydrated votes of a round.
func (service *VotingService) ListVotingRounds(requestContext context.Context) ([]generated.VotingRound, error) {
	slog.DebugContext(requestContext, "service call started")
	votingRoundModels, err := service.votingRoundRepository.ListVotingRounds(requestContext)
	if err != nil {
		return nil, err
	}

	return mapVotingRoundsModelToDTOWithEmptyVotes(votingRoundModels), nil
}

// GetVotingRoundByID returns the full detail of a voting round, including the resolved normative
// and its parent law bucket inline. If the normative or law bucket cannot be resolved, that DTO
// field is left zero-valued and a warning is logged rather than failing the whole request.
func (service *VotingService) GetVotingRoundByID(requestContext context.Context, votingRoundID int) (*generated.VotingRoundDetail, error) {
	slog.DebugContext(requestContext, "service call started", "votingRoundID", votingRoundID)
	votingRoundModel, err := service.votingRoundRepository.GetVotingRoundByID(requestContext, votingRoundID)
	if err != nil {
		return nil, err
	}

	votingRoundDetailDTO := generated.VotingRoundDetail{
		Id:               votingRoundModel.ID,
		Title:            votingRoundModel.Title,
		Description:      votingRoundModel.Description,
		VoteDate:         votingRoundModel.VoteDate,
		NormativeId:      votingRoundModel.NormativeID,
		NormativeVersion: votingRoundModel.NormativeVersion,
		MajorityType:     generated.MajorityType(votingRoundModel.MajorityType),
		Chamber:          generated.Chamber(votingRoundModel.Chamber),
		Votes:            mapVotesModelToDTO(votingRoundModel.Votes),
	}

	normativeModel, err := service.normativeRepository.GetNormativeByIDAndVersion(requestContext, votingRoundModel.NormativeID, votingRoundModel.NormativeVersion)
	if err != nil {
		slog.WarnContext(requestContext, "could not resolve normative for voting round",
			"votingRoundID", votingRoundID,
			"normativeID", votingRoundModel.NormativeID,
			"normativeVersion", votingRoundModel.NormativeVersion,
			"error", err,
		)
		return &votingRoundDetailDTO, nil
	}
	votingRoundDetailDTO.Normative = mapNormativeModelToDTO(*normativeModel)

	lawBucketModel, err := service.lawBucketRepository.GetLawBucketByID(requestContext, normativeModel.LawBucketID)
	if err != nil {
		slog.WarnContext(requestContext, "could not resolve law bucket for voting round",
			"votingRoundID", votingRoundID,
			"lawBucketID", normativeModel.LawBucketID,
			"error", err,
		)
		return &votingRoundDetailDTO, nil
	}
	votingRoundDetailDTO.LawBucket = mapLawBucketModelToDTO(*lawBucketModel)

	return &votingRoundDetailDTO, nil
}

// GetVotesByVotingRound reads the round's embedded votes array directly and hydrates each entry
// with the politician and party it references, preserving the embedded array's order. Politician
// and party lookups are batched (one call each) rather than done per-vote.
func (service *VotingService) GetVotesByVotingRound(requestContext context.Context, votingRoundID int) ([]generated.HydratedVote, error) {
	slog.DebugContext(requestContext, "service call started", "votingRoundID", votingRoundID)
	votingRoundModel, err := service.votingRoundRepository.GetVotingRoundByID(requestContext, votingRoundID)
	if err != nil {
		return nil, err
	}

	politicianIDSet := make(map[string]struct{}, len(votingRoundModel.Votes))
	partyIDSet := make(map[string]struct{}, len(votingRoundModel.Votes))
	for _, voteModel := range votingRoundModel.Votes {
		politicianIDSet[voteModel.PoliticianID] = struct{}{}
		partyIDSet[voteModel.PartyID] = struct{}{}
	}

	politicianIDs := make([]string, 0, len(politicianIDSet))
	for politicianID := range politicianIDSet {
		politicianIDs = append(politicianIDs, politicianID)
	}
	partyIDs := make([]string, 0, len(partyIDSet))
	for partyID := range partyIDSet {
		partyIDs = append(partyIDs, partyID)
	}

	politiciansByID, err := service.politicianRepository.GetPoliticiansByIDs(requestContext, politicianIDs)
	if err != nil {
		return nil, err
	}
	partiesByID, err := service.partyRepository.GetPartiesByIDs(requestContext, partyIDs)
	if err != nil {
		return nil, err
	}

	hydratedVoteDTOs := make([]generated.HydratedVote, 0, len(votingRoundModel.Votes))
	for _, voteModel := range votingRoundModel.Votes {
		hydratedVoteDTO := generated.HydratedVote{
			PoliticianId: voteModel.PoliticianID,
			PartyId:      voteModel.PartyID,
			Value:        generated.VoteValue(voteModel.Value),
		}
		if politicianModel, found := politiciansByID[voteModel.PoliticianID]; found {
			hydratedVoteDTO.Politician = mapPoliticianModelToDTO(politicianModel)
		}
		if partyModel, found := partiesByID[voteModel.PartyID]; found {
			hydratedVoteDTO.Party = mapPartyModelToDTO(partyModel)
		}
		hydratedVoteDTOs = append(hydratedVoteDTOs, hydratedVoteDTO)
	}

	return hydratedVoteDTOs, nil
}
