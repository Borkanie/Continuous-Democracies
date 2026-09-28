package models

import "time"

// LawBucket is the real-world law/bill. Its _id is the natural external key
// (cdep's idp), so no redundant synthetic UUID is stored. Normatives is never
// persisted on this document — it is populated at read time via $lookup,
// embedding the current (highest-version) normative for each distinct
// normative id belonging to this bucket.
type LawBucket struct {
	ID             int         `bson:"_id" json:"id"`
	Version        int         `bson:"version" json:"version"`
	PLNumber       string      `bson:"plNumber" json:"plNumber"`
	Title          string      `bson:"title" json:"title"`
	Description    string      `bson:"description" json:"description"`
	InitiationDate time.Time   `bson:"initiationDate" json:"initiationDate"`
	Status         string      `bson:"status" json:"status"`
	Normatives     []Normative `bson:"normatives,omitempty" json:"normatives"`

	// SenateRegistrationNumber is the Senate's own registration number for
	// this same bill (e.g. "L235/2026"), pairing with PLNumber (the Chamber
	// of Deputies' registration number, e.g. "PLX592/2025") the same way
	// both chambers' own sites cross-reference each other. Empty until the
	// bill has actually reached the Senate / the pairing is known.
	SenateRegistrationNumber string `bson:"senateRegistrationNumber,omitempty" json:"senateRegistrationNumber,omitempty"`
}
