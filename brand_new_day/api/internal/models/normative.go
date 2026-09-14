package models

import "time"

// NormativeKey is the compound _id of a Normative document: Mongo's native
// pattern for versioned documents (no string-concatenation tricks). Two
// Normative documents share the same ID but differ in Version when a law's
// text changes after a vote.
type NormativeKey struct {
	ID      string `bson:"id" json:"id"`
	Version int    `bson:"version" json:"version"`
}

// Normative is one votable unit: an article, amendment, referenced act, or
// the whole-bill final-adoption text. It is a separate collection (not an
// embedded LawBucket.normatives[] array) because VotingRound needs to
// reference one exact version of one normative directly.
type Normative struct {
	Key           NormativeKey `bson:"_id" json:"-"`
	LawBucketID   int          `bson:"lawBucketId" json:"lawBucketId"`
	Type          string       `bson:"type" json:"type"`
	Label         string       `bson:"label" json:"label"`
	Text          string       `bson:"text" json:"text"`
	EffectiveDate time.Time    `bson:"effectiveDate" json:"effectiveDate"`
}
