package service

import (
	"github.com/borkanie/brand-new-day-api/internal/generated"
	"github.com/borkanie/brand-new-day-api/internal/models"
)

// mapPartyModelToDTO converts a models.Party into the generated.Party DTO.
func mapPartyModelToDTO(partyModel models.Party) generated.Party {
	return generated.Party{
		Id:      partyModel.ID,
		Name:    partyModel.Name,
		Acronym: partyModel.Acronym,
		LogoUrl: partyModel.LogoURL,
		Color:   partyModel.Color,
		Active:  partyModel.Active,
	}
}

// mapPoliticianModelToDTO converts a models.Politician into the generated.Politician DTO.
func mapPoliticianModelToDTO(politicianModel models.Politician) generated.Politician {
	return generated.Politician{
		Id:           politicianModel.ID,
		Name:         politicianModel.Name,
		Gender:       politicianModel.Gender,
		ImageUrl:     politicianModel.ImageURL,
		PartyId:      politicianModel.PartyID,
		Active:       politicianModel.Active,
		WorkLocation: politicianModel.WorkLocation,
	}
}

// mapNormativeModelToDTO converts a models.Normative into the generated.Normative DTO.
// Normative has no flat ID/Version fields on the model; they live under Key.
func mapNormativeModelToDTO(normativeModel models.Normative) generated.Normative {
	return generated.Normative{
		Id:            normativeModel.Key.ID,
		Version:       normativeModel.Key.Version,
		LawBucketId:   normativeModel.LawBucketID,
		Type:          generated.NormativeType(normativeModel.Type),
		Label:         normativeModel.Label,
		Text:          normativeModel.Text,
		EffectiveDate: normativeModel.EffectiveDate,
	}
}

// mapNormativesModelToDTO converts a slice of models.Normative into a slice of generated.Normative,
// always returning a non-nil slice so it serializes as [] rather than null.
func mapNormativesModelToDTO(normativeModels []models.Normative) []generated.Normative {
	normativeDTOs := make([]generated.Normative, 0, len(normativeModels))
	for _, normativeModel := range normativeModels {
		normativeDTOs = append(normativeDTOs, mapNormativeModelToDTO(normativeModel))
	}
	return normativeDTOs
}

// mapLawBucketModelToDTO converts a models.LawBucket into the generated.LawBucket DTO,
// mapping the already-embedded current normatives through as well.
func mapLawBucketModelToDTO(lawBucketModel models.LawBucket) generated.LawBucket {
	return generated.LawBucket{
		Id:             lawBucketModel.ID,
		Version:        lawBucketModel.Version,
		PlNumber:       lawBucketModel.PLNumber,
		Title:          lawBucketModel.Title,
		Description:    lawBucketModel.Description,
		InitiationDate: lawBucketModel.InitiationDate,
		Status:         lawBucketModel.Status,
		Normatives:     mapNormativesModelToDTO(lawBucketModel.Normatives),
	}
}

// mapVoteModelToDTO converts a models.Vote into the generated.Vote DTO.
func mapVoteModelToDTO(voteModel models.Vote) generated.Vote {
	return generated.Vote{
		PoliticianId: voteModel.PoliticianID,
		PartyId:      voteModel.PartyID,
		Value:        generated.VoteValue(voteModel.Value),
	}
}

// mapVotesModelToDTO converts a slice of models.Vote into a slice of generated.Vote,
// always returning a non-nil slice so it serializes as [] rather than null.
func mapVotesModelToDTO(voteModels []models.Vote) []generated.Vote {
	voteDTOs := make([]generated.Vote, 0, len(voteModels))
	for _, voteModel := range voteModels {
		voteDTOs = append(voteDTOs, mapVoteModelToDTO(voteModel))
	}
	return voteDTOs
}

// mapVotingRoundModelToDTO converts a models.VotingRound into the generated.VotingRound DTO,
// including its full embedded votes array.
func mapVotingRoundModelToDTO(votingRoundModel models.VotingRound) generated.VotingRound {
	return generated.VotingRound{
		Id:               votingRoundModel.ID,
		Title:            votingRoundModel.Title,
		Description:      votingRoundModel.Description,
		VoteDate:         votingRoundModel.VoteDate,
		NormativeId:      votingRoundModel.NormativeID,
		NormativeVersion: votingRoundModel.NormativeVersion,
		MajorityType:     generated.MajorityType(votingRoundModel.MajorityType),
		Votes:            mapVotesModelToDTO(votingRoundModel.Votes),
	}
}

// mapVotingRoundModelToDTOWithEmptyVotes converts a models.VotingRound into the generated.VotingRound
// DTO but always sets Votes to an empty slice. Used for list-style endpoints (ListVotingRounds,
// GetVotingRoundsByLawBucket) so the response doesn't inline hundreds of votes per round; callers
// should use GET /votingRounds/{roundId}/votes to fetch the hydrated votes for a specific round.
func mapVotingRoundModelToDTOWithEmptyVotes(votingRoundModel models.VotingRound) generated.VotingRound {
	votingRoundDTO := mapVotingRoundModelToDTO(votingRoundModel)
	votingRoundDTO.Votes = []generated.Vote{}
	return votingRoundDTO
}

// mapVotingRoundsModelToDTOWithEmptyVotes maps a slice of models.VotingRound using
// mapVotingRoundModelToDTOWithEmptyVotes, always returning a non-nil slice.
func mapVotingRoundsModelToDTOWithEmptyVotes(votingRoundModels []models.VotingRound) []generated.VotingRound {
	votingRoundDTOs := make([]generated.VotingRound, 0, len(votingRoundModels))
	for _, votingRoundModel := range votingRoundModels {
		votingRoundDTOs = append(votingRoundDTOs, mapVotingRoundModelToDTOWithEmptyVotes(votingRoundModel))
	}
	return votingRoundDTOs
}
