package models

// Vote is the DB document stored in MongoDB.
type Vote struct {
	ID           string `bson:"_id"          json:"id"`
	PoliticianId string `bson:"politicianId" json:"-"`
	RoundId      string `bson:"roundId"      json:"-"`
	Position     int    `bson:"position"     json:"position"`
}

// VoteResult is returned by aggregation queries that join politician, party, and round.
type VoteResult struct {
	ID         string               `bson:"_id"        json:"id"`
	Position   int                  `bson:"position"   json:"position"`
	Politician *PoliticianWithParty `bson:"politician" json:"politician"`
	Round      *VotingRound         `bson:"round"      json:"round"`
	Name       string               `bson:"name"       json:"name"` // = politician.Name
}
