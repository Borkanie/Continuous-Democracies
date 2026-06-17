package models

import "time"

// VotingRound is the DB document. Name is kept equal to Title for frontend compatibility.
type VotingRound struct {
	ID          string    `bson:"_id"         json:"id"`
	VoteId      int       `bson:"voteId"      json:"voteId"`
	Title       string    `bson:"title"       json:"title"`
	Name        string    `bson:"name"        json:"name"`
	Description string    `bson:"description" json:"description"`
	VoteDate    time.Time `bson:"voteDate"    json:"voteDate"`
}
