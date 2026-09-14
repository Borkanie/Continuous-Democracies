// Command seed populates a fresh MongoDB database with the mock dataset
// described in brand_new_day/ARCHITECTURE.md's "Seed data" section: 5
// parties, 330 politicians distributed across them, 10 law buckets each with
// 1-3 normatives (a mix of article/amendment/referencedAct/wholeBillFinal,
// with at least one normative carrying 2 versions) and one voting round per
// normative version, with every politician casting a randomly-assigned vote
// on every round.
//
// Politicians and votes are generated in code with a fixed RNG seed so that
// `make seed` against a fresh database is fully reproducible. Parties and law
// buckets/normatives/voting-round metadata come from the hand-written JSON
// fixtures under fixtures/, embedded at build time.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"time"

	"github.com/borkanie/brand-new-day-api/fixtures"
	"github.com/borkanie/brand-new-day-api/internal/config"
	"github.com/borkanie/brand-new-day-api/internal/db"
	"github.com/borkanie/brand-new-day-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// seedRNGSeed1/seedRNGSeed2 fix math/rand/v2's PCG source so re-running the
// seeder against a fresh database always produces byte-for-byte the same
// politicians, name assignments and votes.
const (
	seedRNGSeed1 = 42
	seedRNGSeed2 = 42

	politicianCount = 330
	dateLayout      = "2006-01-02"
	seedTimeout     = 60 * time.Second
)

// --- fixture JSON shapes -----------------------------------------------

type partyFixture struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
	LogoURL string `json:"logoUrl"`
	Color   string `json:"color"`
	Active  bool   `json:"active"`
}

type votingRoundFixture struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	VoteDate     string `json:"voteDate"`
	MajorityType string `json:"majorityType"`
}

type normativeVersionFixture struct {
	Version       int                `json:"version"`
	Type          string             `json:"type"`
	Label         string             `json:"label"`
	Text          string             `json:"text"`
	EffectiveDate string             `json:"effectiveDate"`
	VotingRound   votingRoundFixture `json:"votingRound"`
}

type normativeFixture struct {
	ID       string                    `json:"id"`
	Versions []normativeVersionFixture `json:"versions"`
}

type lawBucketFixture struct {
	ID             int                `json:"id"`
	Version        int                `json:"version"`
	PLNumber       string             `json:"plNumber"`
	Title          string             `json:"title"`
	Description    string             `json:"description"`
	InitiationDate string             `json:"initiationDate"`
	Status         string             `json:"status"`
	Normatives     []normativeFixture `json:"normatives"`
}

// --- Romanian name pools for deterministic politician generation --------

var maleFirstNames = []string{
	"Andrei", "Ion", "Mihai", "Cristian", "Vasile", "Gheorghe", "Alexandru",
	"Florin", "Marian", "Dan", "Bogdan", "Radu", "Sorin", "Adrian",
	"Cătălin", "Nicolae", "Constantin", "Ștefan", "Valentin", "Daniel",
}

var femaleFirstNames = []string{
	"Maria", "Elena", "Ioana", "Ana", "Cristina", "Gabriela", "Andreea",
	"Mihaela", "Simona", "Raluca", "Diana", "Alina", "Camelia",
	"Georgiana", "Luminița", "Carmen", "Roxana", "Daniela", "Monica", "Corina",
}

var lastNames = []string{
	"Popescu", "Ionescu", "Popa", "Radu", "Dumitru", "Stoica", "Gheorghiu",
	"Constantin", "Marin", "Tudor", "Matei", "Stan", "Neagu", "Barbu",
	"Florea", "Nistor", "Enache", "Voicu", "Iliescu", "Cristea", "Vasilescu",
	"Manea", "Toma", "Dobre", "Preda", "Nedelcu", "Anghel", "Luca",
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	configuration := config.Load()

	seedContext, cancelSeedContext := context.WithTimeout(context.Background(), seedTimeout)
	defer cancelSeedContext()

	mongoClient, database, err := db.ConnectClient(configuration.MongoURI, configuration.DatabaseName)
	if err != nil {
		slog.Error("failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	defer func() {
		disconnectContext, cancelDisconnect := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelDisconnect()
		if err := mongoClient.Disconnect(disconnectContext); err != nil {
			slog.Warn("failed to disconnect MongoDB client", "error", err)
		}
	}()

	if err := dropSeedCollections(seedContext, database); err != nil {
		slog.Error("failed to drop existing collections", "error", err)
		os.Exit(1)
	}

	parties, err := loadPartyFixtures()
	if err != nil {
		slog.Error("failed to load party fixtures", "error", err)
		os.Exit(1)
	}

	lawBuckets, err := loadLawBucketFixtures()
	if err != nil {
		slog.Error("failed to load law bucket fixtures", "error", err)
		os.Exit(1)
	}

	randomNumberGenerator := rand.New(rand.NewPCG(seedRNGSeed1, seedRNGSeed2))

	politicians := generatePoliticians(randomNumberGenerator, parties)

	insertedLawBucketCount, insertedNormativeCount, insertedVotingRoundCount, insertedVoteCount, err := seedLawData(
		seedContext, database, randomNumberGenerator, lawBuckets, politicians,
	)
	if err != nil {
		slog.Error("failed to seed law data", "error", err)
		os.Exit(1)
	}

	if err := insertParties(seedContext, database, parties); err != nil {
		slog.Error("failed to insert parties", "error", err)
		os.Exit(1)
	}

	if err := insertPoliticians(seedContext, database, politicians); err != nil {
		slog.Error("failed to insert politicians", "error", err)
		os.Exit(1)
	}

	if err := db.EnsureIndexes(seedContext, database); err != nil {
		slog.Error("failed to ensure indexes", "error", err)
		os.Exit(1)
	}

	slog.Info("seed complete",
		"parties", len(parties),
		"politicians", len(politicians),
		"lawBuckets", insertedLawBucketCount,
		"normatives", insertedNormativeCount,
		"votingRounds", insertedVotingRoundCount,
		"votes", insertedVoteCount,
	)
}

// dropSeedCollections drops every collection this seeder owns, so re-running
// `make seed` against an already-seeded database produces a clean, identical
// dataset rather than accumulating duplicates.
func dropSeedCollections(requestContext context.Context, database *mongo.Database) error {
	collectionNames := []string{
		db.CollectionParties,
		db.CollectionPoliticians,
		db.CollectionLawBuckets,
		db.CollectionNormatives,
		db.CollectionVotingRounds,
	}
	for _, collectionName := range collectionNames {
		if err := database.Collection(collectionName).Drop(requestContext); err != nil {
			return fmt.Errorf("drop %s: %w", collectionName, err)
		}
	}
	return nil
}

func loadPartyFixtures() ([]partyFixture, error) {
	rawParties, err := fixtures.FS.ReadFile("parties.json")
	if err != nil {
		return nil, fmt.Errorf("read parties.json: %w", err)
	}
	var parties []partyFixture
	if err := json.Unmarshal(rawParties, &parties); err != nil {
		return nil, fmt.Errorf("unmarshal parties.json: %w", err)
	}
	return parties, nil
}

func loadLawBucketFixtures() ([]lawBucketFixture, error) {
	rawLawBuckets, err := fixtures.FS.ReadFile("law_buckets.json")
	if err != nil {
		return nil, fmt.Errorf("read law_buckets.json: %w", err)
	}
	var lawBuckets []lawBucketFixture
	if err := json.Unmarshal(rawLawBuckets, &lawBuckets); err != nil {
		return nil, fmt.Errorf("unmarshal law_buckets.json: %w", err)
	}
	return lawBuckets, nil
}

// generatePoliticians deterministically builds the 330-member steady-state
// parliament, distributed evenly (66 each) round-robin across the 5 seeded
// parties.
func generatePoliticians(randomNumberGenerator *rand.Rand, parties []partyFixture) []models.Politician {
	politicians := make([]models.Politician, 0, politicianCount)

	for politicianIndex := 0; politicianIndex < politicianCount; politicianIndex++ {
		party := parties[politicianIndex%len(parties)]

		politicianID := fmt.Sprintf("pol-%03d", politicianIndex+1)

		var gender int
		var firstName string
		if randomNumberGenerator.IntN(2) == 0 {
			gender = 0
			firstName = maleFirstNames[randomNumberGenerator.IntN(len(maleFirstNames))]
		} else {
			gender = 1
			firstName = femaleFirstNames[randomNumberGenerator.IntN(len(femaleFirstNames))]
		}
		lastName := lastNames[randomNumberGenerator.IntN(len(lastNames))]

		// ~95% of the seeded parliament is active, matching a realistic
		// steady-state mix of a handful of vacant/suspended mandates.
		active := randomNumberGenerator.Float64() < 0.95

		// Romania has 41 counties plus Bucharest municipality: 42 plausible
		// work-location codes.
		workLocation := randomNumberGenerator.IntN(42) + 1

		politicians = append(politicians, models.Politician{
			ID:           politicianID,
			Name:         fmt.Sprintf("%s %s", firstName, lastName),
			Gender:       gender,
			ImageURL:     fmt.Sprintf("https://example.com/politicians/%s.jpg", politicianID),
			PartyID:      party.ID,
			Active:       active,
			WorkLocation: workLocation,
		})
	}

	return politicians
}

// voteValuePool returns the weighted pool of vote values a round's votes are
// drawn from. Simple-majority rounds keep the original uniform distribution
// (Yes>No is a coin flip either way). Absolute/qualified rounds need a much
// higher Yes count to pass (>165 / >=220 out of 330) than a uniform draw ever
// produces, so their pools are biased toward Yes - still randomized off the
// same fixed seed, just giving the demo data a plausible shot at showing a
// PASSED absolute/qualified round instead of guaranteeing FAILED every time.
func voteValuePool(majorityType string) []string {
	switch majorityType {
	case models.MajorityTypeQualified:
		return []string{
			models.VoteValueYes, models.VoteValueYes, models.VoteValueYes,
			models.VoteValueYes, models.VoteValueYes, models.VoteValueYes,
			models.VoteValueYes, models.VoteValueYes, models.VoteValueYes,
			models.VoteValueYes, models.VoteValueYes, models.VoteValueYes,
			models.VoteValueYes,
			models.VoteValueNo, models.VoteValueNo,
			models.VoteValueAbstain, models.VoteValueAbsent,
		}
	case models.MajorityTypeAbsolute:
		return []string{
			models.VoteValueYes, models.VoteValueYes, models.VoteValueYes,
			models.VoteValueYes, models.VoteValueYes, models.VoteValueYes,
			models.VoteValueNo, models.VoteValueNo, models.VoteValueNo,
			models.VoteValueNo,
			models.VoteValueAbstain, models.VoteValueAbsent,
		}
	default:
		return models.AllVoteValues()
	}
}

// seedLawData walks the law bucket fixtures, builds and inserts the
// LawBucket/Normative/VotingRound documents (with every politician's vote
// embedded on every round), and returns the counts inserted.
func seedLawData(
	requestContext context.Context,
	database *mongo.Database,
	randomNumberGenerator *rand.Rand,
	lawBucketFixtures []lawBucketFixture,
	politicians []models.Politician,
) (lawBucketCount int, normativeCount int, votingRoundCount int, voteCount int, err error) {
	lawBucketsCollection := database.Collection(db.CollectionLawBuckets)
	normativesCollection := database.Collection(db.CollectionNormatives)
	votingRoundsCollection := database.Collection(db.CollectionVotingRounds)

	var lawBucketDocuments []any
	var normativeDocuments []any
	var votingRoundDocuments []any

	for _, lawBucket := range lawBucketFixtures {
		initiationDate, parseErr := time.Parse(dateLayout, lawBucket.InitiationDate)
		if parseErr != nil {
			return 0, 0, 0, 0, fmt.Errorf("parse initiationDate for law bucket %d: %w", lawBucket.ID, parseErr)
		}

		lawBucketDocuments = append(lawBucketDocuments, models.LawBucket{
			ID:             lawBucket.ID,
			Version:        lawBucket.Version,
			PLNumber:       lawBucket.PLNumber,
			Title:          lawBucket.Title,
			Description:    lawBucket.Description,
			InitiationDate: initiationDate,
			Status:         lawBucket.Status,
		})

		for _, normative := range lawBucket.Normatives {
			for _, normativeVersion := range normative.Versions {
				effectiveDate, parseErr := time.Parse(dateLayout, normativeVersion.EffectiveDate)
				if parseErr != nil {
					return 0, 0, 0, 0, fmt.Errorf("parse effectiveDate for normative %s v%d: %w", normative.ID, normativeVersion.Version, parseErr)
				}

				normativeDocuments = append(normativeDocuments, models.Normative{
					Key: models.NormativeKey{
						ID:      normative.ID,
						Version: normativeVersion.Version,
					},
					LawBucketID:   lawBucket.ID,
					Type:          normativeVersion.Type,
					Label:         normativeVersion.Label,
					Text:          normativeVersion.Text,
					EffectiveDate: effectiveDate,
				})

				voteDate, parseErr := time.Parse(dateLayout, normativeVersion.VotingRound.VoteDate)
				if parseErr != nil {
					return 0, 0, 0, 0, fmt.Errorf("parse voteDate for voting round %d: %w", normativeVersion.VotingRound.ID, parseErr)
				}

				votePool := voteValuePool(normativeVersion.VotingRound.MajorityType)
				votes := make([]models.Vote, 0, len(politicians))
				for _, politician := range politicians {
					voteValue := votePool[randomNumberGenerator.IntN(len(votePool))]
					votes = append(votes, models.Vote{
						PoliticianID: politician.ID,
						PartyID:      politician.PartyID,
						Value:        voteValue,
					})
				}
				voteCount += len(votes)

				votingRoundDocuments = append(votingRoundDocuments, models.VotingRound{
					ID:               normativeVersion.VotingRound.ID,
					Title:            normativeVersion.VotingRound.Title,
					Description:      normativeVersion.VotingRound.Description,
					VoteDate:         voteDate,
					NormativeID:      normative.ID,
					NormativeVersion: normativeVersion.Version,
					MajorityType:     normativeVersion.VotingRound.MajorityType,
					Votes:            votes,
				})
				votingRoundCount++
			}
			normativeCount += len(normative.Versions)
		}
		lawBucketCount++
	}

	if _, err := lawBucketsCollection.InsertMany(requestContext, lawBucketDocuments); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("insert law buckets: %w", err)
	}
	if _, err := normativesCollection.InsertMany(requestContext, normativeDocuments); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("insert normatives: %w", err)
	}
	if _, err := votingRoundsCollection.InsertMany(requestContext, votingRoundDocuments); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("insert voting rounds: %w", err)
	}

	return lawBucketCount, normativeCount, votingRoundCount, voteCount, nil
}

func insertParties(requestContext context.Context, database *mongo.Database, partyFixtures []partyFixture) error {
	partyDocuments := make([]any, 0, len(partyFixtures))
	for _, party := range partyFixtures {
		partyDocuments = append(partyDocuments, models.Party{
			ID:      party.ID,
			Name:    party.Name,
			Acronym: party.Acronym,
			LogoURL: party.LogoURL,
			Color:   party.Color,
			Active:  party.Active,
		})
	}
	if _, err := database.Collection(db.CollectionParties).InsertMany(requestContext, partyDocuments); err != nil {
		return fmt.Errorf("insert parties: %w", err)
	}
	return nil
}

func insertPoliticians(requestContext context.Context, database *mongo.Database, politicians []models.Politician) error {
	politicianDocuments := make([]any, 0, len(politicians))
	for _, politician := range politicians {
		politicianDocuments = append(politicianDocuments, politician)
	}
	if _, err := database.Collection(db.CollectionPoliticians).InsertMany(requestContext, politicianDocuments); err != nil {
		return fmt.Errorf("insert politicians: %w", err)
	}
	return nil
}
