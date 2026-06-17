package models

// Politician is the DB document stored in MongoDB.
type Politician struct {
	ID           string  `bson:"_id"          json:"id"`
	Name         string  `bson:"name"         json:"name"`
	Gender       int     `bson:"gender"       json:"gender"`
	ImageUrl     *string `bson:"imageUrl"     json:"imageUrl"`
	PartyId      string  `bson:"partyId"      json:"-"`
	Active       bool    `bson:"active"       json:"active"`
	WorkLocation int     `bson:"workLocation" json:"workLocation"`
}

// PoliticianWithParty is returned by aggregation queries that join to the parties collection.
type PoliticianWithParty struct {
	ID           string  `bson:"_id"          json:"id"`
	Name         string  `bson:"name"         json:"name"`
	Gender       int     `bson:"gender"       json:"gender"`
	ImageUrl     *string `bson:"imageUrl"     json:"imageUrl"`
	PartyId      string  `bson:"partyId"      json:"-"`
	Active       bool    `bson:"active"       json:"active"`
	WorkLocation int     `bson:"workLocation" json:"workLocation"`
	Party        *Party  `bson:"party"        json:"party"`
}
