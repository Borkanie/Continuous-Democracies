package models

import "time"

type Party struct {
	ID      string  `bson:"_id"`
	Name    string  `bson:"name"`
	Acronym string  `bson:"acronym"`
	LogoUrl *string `bson:"logoUrl"`
	Color   string  `bson:"color"`
	Active  bool    `bson:"active"`
}

type Politician struct {
	ID           string  `bson:"_id"`
	Name         string  `bson:"name"`
	Gender       int     `bson:"gender"`
	ImageUrl     *string `bson:"imageUrl"`
	PartyId      string  `bson:"partyId"`
	Active       bool    `bson:"active"`
	WorkLocation int     `bson:"workLocation"`
}

type VotingRound struct {
	ID          string    `bson:"_id"`
	VoteId      int       `bson:"voteId"`
	Title       string    `bson:"title"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	VoteDate    time.Time `bson:"voteDate"`
}

type Vote struct {
	ID           string `bson:"_id"`
	PoliticianId string `bson:"politicianId"`
	RoundId      string `bson:"roundId"`
	Position     int    `bson:"position"`
}

// VotePosition constants mirror the Romanian parliament vote values
const (
	PositionYes     = 0
	PositionNo      = 1
	PositionAbstain = 2
	PositionAbsent  = 3
)
