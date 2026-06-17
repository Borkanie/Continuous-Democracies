package models

type Party struct {
	ID      string  `bson:"_id"    json:"id"`
	Name    string  `bson:"name"   json:"name"`
	Acronym string  `bson:"acronym" json:"acronym"`
	LogoUrl *string `bson:"logoUrl" json:"logoUrl"`
	Color   string  `bson:"color"  json:"color"`
	Active  bool    `bson:"active" json:"active"`
}
