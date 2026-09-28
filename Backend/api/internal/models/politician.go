package models

// Politician is a member of parliament. It keeps a synthetic string
// identifier because no natural external key exists for politicians in the
// source data.
type Politician struct {
	ID           string `bson:"_id" json:"id"`
	Name         string `bson:"name" json:"name"`
	Gender       int    `bson:"gender" json:"gender"`
	ImageURL     string `bson:"imageUrl" json:"imageUrl"`
	PartyID      string `bson:"partyId" json:"partyId"`
	Active       bool   `bson:"active" json:"active"`
	WorkLocation int    `bson:"workLocation" json:"workLocation"`
}
