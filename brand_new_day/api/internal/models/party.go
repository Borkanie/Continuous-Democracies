package models

// Party is a political party. It keeps a synthetic string identifier because
// no natural external key exists for parties in the source data.
type Party struct {
	ID      string `bson:"_id" json:"id"`
	Name    string `bson:"name" json:"name"`
	Acronym string `bson:"acronym" json:"acronym"`
	LogoURL string `bson:"logoUrl" json:"logoUrl"`
	Color   string `bson:"color" json:"color"`
	Active  bool   `bson:"active" json:"active"`
}
