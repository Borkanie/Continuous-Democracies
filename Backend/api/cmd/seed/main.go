// Command seed populates a fresh MongoDB database with mostly-real data: the
// real current Senate roster (senators.json, scraped from
// senat.ro/FisaSenatori.aspx - real names, real parliamentary group, real
// electoral circumscription) and the real current Chamber of Deputies
// roster (deputies.json, scraped from cdep.ro's structura2015.de pages,
// same real fields), plus the parties/political groups those rosters
// reference (parties.json). Law buckets/normatives/voting-round metadata
// (fixtures/law_buckets.json) mix real historical bills (with real dates and
// real recorded vote tallies where researched - see each entry's
// description) with invented demo bills; every voting round's embedded
// votes are either a real tally (VotingRound.Tally) or a reproducible random
// draw, cast by whichever real roster matches that round's Chamber - Senate
// rounds draw from senators.json, Camera rounds from deputies.json.
//
// What's NOT real: individual vote choices on invented/non-tallied rounds
// (randomly assigned), and politician photos (senat.ro/cdep.ro's roster
// pages don't expose them without per-person scraping, so ImageURL is still
// a placeholder). Random draws use a fixed RNG seed so `make seed` against a
// fresh database is fully reproducible.
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
// non-tallied vote assignments.
const (
	seedRNGSeed1 = 42
	seedRNGSeed2 = 42

	dateLayout  = "2006-01-02"
	seedTimeout = 60 * time.Second
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
	ID           int               `json:"id"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	VoteDate     string            `json:"voteDate"`
	MajorityType string            `json:"majorityType"`
	Chamber      string            `json:"chamber"`
	Tally        *voteTallyFixture `json:"tally,omitempty"`
}

// voteTallyFixture pins a round's Yes/No/Abstain counts to a real recorded
// result instead of letting the random pool decide. Any politician in that
// round's real roster (senators.json for Chamber=senate, deputies.json for
// Chamber=parliament) not covered by Yes+No+Abstain is marked Absent. The
// scraped rosters are a few seats short of the nominal 136/330 (a handful of
// vacant seats or parsing misses), so this Absent count is close to but not
// always exactly the real session's reported absentee figure.
type voteTallyFixture struct {
	Yes     int `json:"yes"`
	No      int `json:"no"`
	Abstain int `json:"abstain"`
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
	ID                       int                `json:"id"`
	Version                  int                `json:"version"`
	PLNumber                 string             `json:"plNumber"`
	SenateRegistrationNumber string             `json:"senateRegistrationNumber"`
	Title                    string             `json:"title"`
	Description              string             `json:"description"`
	InitiationDate           string             `json:"initiationDate"`
	Status                   string             `json:"status"`
	Normatives               []normativeFixture `json:"normatives"`
}

// politicianFixture mirrors senators.json/deputies.json - real roster
// entries scraped from senat.ro/cdep.ro, already shaped like
// models.Politician (see loadPoliticianFixtures).
type politicianFixture struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Gender       int    `json:"gender"`
	ImageURL     string `json:"imageUrl"`
	PartyID      string `json:"partyId"`
	Active       bool   `json:"active"`
	WorkLocation int    `json:"workLocation"`
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

	senators, err := loadPoliticianFixtures("senators.json")
	if err != nil {
		slog.Error("failed to load senator fixtures", "error", err)
		os.Exit(1)
	}

	deputies, err := loadPoliticianFixtures("deputies.json")
	if err != nil {
		slog.Error("failed to load deputy fixtures", "error", err)
		os.Exit(1)
	}

	randomNumberGenerator := rand.New(rand.NewPCG(seedRNGSeed1, seedRNGSeed2))

	rosters := chamberRosters{senate: senators, parliament: deputies}

	insertedLawBucketCount, insertedNormativeCount, insertedVotingRoundCount, insertedVoteCount, err := seedLawData(
		seedContext, database, randomNumberGenerator, lawBuckets, rosters,
	)
	if err != nil {
		slog.Error("failed to seed law data", "error", err)
		os.Exit(1)
	}

	if err := insertParties(seedContext, database, parties); err != nil {
		slog.Error("failed to insert parties", "error", err)
		os.Exit(1)
	}

	allPoliticians := append(append([]models.Politician{}, senators...), deputies...)
	if err := insertPoliticians(seedContext, database, allPoliticians); err != nil {
		slog.Error("failed to insert politicians", "error", err)
		os.Exit(1)
	}

	if err := db.EnsureIndexes(seedContext, database); err != nil {
		slog.Error("failed to ensure indexes", "error", err)
		os.Exit(1)
	}

	slog.Info("seed complete",
		"parties", len(parties),
		"senators", len(senators),
		"deputies", len(deputies),
		"lawBuckets", insertedLawBucketCount,
		"normatives", insertedNormativeCount,
		"votingRounds", insertedVotingRoundCount,
		"votes", insertedVoteCount,
	)
}

// chamberRosters holds the two real, disjoint politician pools a voting
// round draws its votes from, keyed by the same Chamber values the round
// fixtures use (models.ChamberSenate / models.ChamberParliament).
type chamberRosters struct {
	senate     []models.Politician
	parliament []models.Politician
}

func (r chamberRosters) forChamber(chamber string) []models.Politician {
	if chamber == models.ChamberSenate {
		return r.senate
	}
	return r.parliament
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

// loadPoliticianFixtures loads a real roster (senators.json or
// deputies.json) straight into models.Politician - the fixture shape
// already matches the model field-for-field, no synthesis needed.
func loadPoliticianFixtures(filename string) ([]models.Politician, error) {
	raw, err := fixtures.FS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filename, err)
	}
	var fixturePoliticians []politicianFixture
	if err := json.Unmarshal(raw, &fixturePoliticians); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", filename, err)
	}

	politicians := make([]models.Politician, 0, len(fixturePoliticians))
	for _, p := range fixturePoliticians {
		politicians = append(politicians, models.Politician{
			ID:           p.ID,
			Name:         p.Name,
			Gender:       p.Gender,
			ImageURL:     p.ImageURL,
			PartyID:      p.PartyID,
			Active:       p.Active,
			WorkLocation: p.WorkLocation,
		})
	}
	return politicians, nil
}

// voteValuePool returns the weighted pool of vote values a round's votes are
// drawn from. Simple-majority rounds keep the original uniform distribution
// (Yes>No is a coin flip either way). Absolute/qualified rounds need a much
// higher Yes count to pass than a uniform draw ever produces, so their pools
// are biased toward Yes - still randomized off the same fixed seed, just
// giving the demo data a plausible shot at showing a PASSED absolute/
// qualified round instead of guaranteeing FAILED every time. The bias ratio
// is independent of chamber size, so it works for both the Senate's and the
// Chamber's rosters.
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

// buildTalliedVotes deterministically assigns a real recorded Yes/No/Abstain
// result across the seeded politician pool (shuffled via the fixed-seed RNG
// so the result is reproducible but not just the first N politicians by ID),
// with every remaining politician marked Absent.
func buildTalliedVotes(
	randomNumberGenerator *rand.Rand,
	politicians []models.Politician,
	tally *voteTallyFixture,
) []models.Vote {
	order := randomNumberGenerator.Perm(len(politicians))
	votes := make([]models.Vote, len(politicians))

	for position, politicianIndex := range order {
		politician := politicians[politicianIndex]

		var voteValue string
		switch {
		case position < tally.Yes:
			voteValue = models.VoteValueYes
		case position < tally.Yes+tally.No:
			voteValue = models.VoteValueNo
		case position < tally.Yes+tally.No+tally.Abstain:
			voteValue = models.VoteValueAbstain
		default:
			voteValue = models.VoteValueAbsent
		}

		votes[politicianIndex] = models.Vote{
			PoliticianID: politician.ID,
			PartyID:      politician.PartyID,
			Value:        voteValue,
		}
	}

	return votes
}

// seedLawData walks the law bucket fixtures, builds and inserts the
// LawBucket/Normative/VotingRound documents (with every politician's vote
// embedded on every round), and returns the counts inserted.
func seedLawData(
	requestContext context.Context,
	database *mongo.Database,
	randomNumberGenerator *rand.Rand,
	lawBucketFixtures []lawBucketFixture,
	rosters chamberRosters,
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
			ID:                       lawBucket.ID,
			Version:                  lawBucket.Version,
			PLNumber:                 lawBucket.PLNumber,
			SenateRegistrationNumber: lawBucket.SenateRegistrationNumber,
			Title:                    lawBucket.Title,
			Description:              lawBucket.Description,
			InitiationDate:           initiationDate,
			Status:                   lawBucket.Status,
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

				politicians := rosters.forChamber(normativeVersion.VotingRound.Chamber)

				var votes []models.Vote
				if tally := normativeVersion.VotingRound.Tally; tally != nil {
					votes = buildTalliedVotes(randomNumberGenerator, politicians, tally)
				} else {
					votePool := voteValuePool(normativeVersion.VotingRound.MajorityType)
					votes = make([]models.Vote, 0, len(politicians))
					for _, politician := range politicians {
						voteValue := votePool[randomNumberGenerator.IntN(len(votePool))]
						votes = append(votes, models.Vote{
							PoliticianID: politician.ID,
							PartyID:      politician.PartyID,
							Value:        voteValue,
						})
					}
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
					Chamber:          normativeVersion.VotingRound.Chamber,
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
