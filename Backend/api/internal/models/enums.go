// Package models defines the persistence-level data shapes for brand_new_day/api.
// Field names and bson tags are the cross-worker contract shared with the
// openapi.yaml / generated types owned by the other workers — do not rename.
package models

// Vote value enum. Exact casing matters: it is serialized as-is into JSON
// responses and matches the openapi.yaml enum.
const (
	VoteValueYes     = "Yes"
	VoteValueNo      = "No"
	VoteValueAbstain = "Abstain"
	VoteValueAbsent  = "Absent"
)

// AllVoteValues returns every valid Vote.Value in a stable order.
func AllVoteValues() []string {
	return []string{VoteValueYes, VoteValueNo, VoteValueAbstain, VoteValueAbsent}
}

// Normative type enum. Exact casing matters: it is serialized as-is into
// JSON responses and matches the openapi.yaml enum.
const (
	NormativeTypeArticle        = "article"
	NormativeTypeAmendment      = "amendment"
	NormativeTypeReferencedAct  = "referencedAct"
	NormativeTypeWholeBillFinal = "wholeBillFinal"
)

// AllNormativeTypes returns every valid Normative.Type in a stable order.
func AllNormativeTypes() []string {
	return []string{
		NormativeTypeArticle,
		NormativeTypeAmendment,
		NormativeTypeReferencedAct,
		NormativeTypeWholeBillFinal,
	}
}

// Majority type enum: which threshold a VotingRound must clear to pass.
// Exact casing matters: it is serialized as-is into JSON responses and
// matches the openapi.yaml enum.
const (
	MajorityTypeSimple    = "simple"
	MajorityTypeAbsolute  = "absolute"
	MajorityTypeQualified = "qualified"
)

// AllMajorityTypes returns every valid VotingRound.MajorityType in a stable order.
func AllMajorityTypes() []string {
	return []string{MajorityTypeSimple, MajorityTypeAbsolute, MajorityTypeQualified}
}

// Chamber enum: which house of Parliament a VotingRound was cast in. Exact
// casing matters: it is serialized as-is into JSON responses and matches the
// openapi.yaml enum.
const (
	ChamberParliament = "parliament"
	ChamberSenate     = "senate"
)

// AllChambers returns every valid VotingRound.Chamber in a stable order.
func AllChambers() []string {
	return []string{ChamberParliament, ChamberSenate}
}
