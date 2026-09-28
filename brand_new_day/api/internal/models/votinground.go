package models

import "time"

// Vote is a single politician's position in a VotingRound. It is not its own
// collection: votes cannot outlive their round and the array is bounded
// (<= 330 politicians), so it is embedded directly in VotingRound.Votes.
// PartyID is snapshotted at vote time so a later party switch doesn't rewrite
// history.
type Vote struct {
	PoliticianID string `bson:"politicianId" json:"politicianId"`
	PartyID      string `bson:"partyId" json:"partyId"`
	Value        string `bson:"value" json:"value"`
}

// VotingRound is one atomic vote event. Its _id is the natural external key
// (cdep's idv), so no redundant synthetic UUID is stored. NormativeID and
// NormativeVersion together are the FK pair into Normative._id, resolving to
// "the law" via Normative -> LawBucket whether this is an article-level vote
// or a whole-bill final vote.
type VotingRound struct {
	ID               int       `bson:"_id" json:"id"`
	Title            string    `bson:"title" json:"title"`
	Description      string    `bson:"description" json:"description"`
	VoteDate         time.Time `bson:"voteDate" json:"voteDate"`
	NormativeID      string    `bson:"normativeId" json:"normativeId"`
	NormativeVersion int       `bson:"normativeVersion" json:"normativeVersion"`
	MajorityType     string    `bson:"majorityType" json:"majorityType"`
	Chamber          string    `bson:"chamber" json:"chamber"`
	Votes            []Vote    `bson:"votes" json:"votes"`
}

// PoliticianVoteRecord is the output shape of the "all votes by this
// politician" aggregation over voting_rounds.votes. VotingRound.Votes is
// always an empty slice in this shape (projected as $literal: []) so a
// politician's vote history isn't ballooned by every other politician's vote
// on the same round.
type PoliticianVoteRecord struct {
	Value       string      `bson:"value"`
	VotingRound VotingRound `bson:"votingRound"`
	Normative   *Normative  `bson:"normative"`
	LawBucket   *LawBucket  `bson:"lawBucket"`
}
