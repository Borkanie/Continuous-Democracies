package poc_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
	"context"

	"github.com/borkanie/parliament-api/internal/db"
	"github.com/borkanie/parliament-api/internal/handlers"
	"github.com/borkanie/parliament-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	pocServer *httptest.Server
	lawId     int
)

func TestMain(m *testing.M) {
	mongoURI := os.Getenv("TEST_MONGO_URI")
	if mongoURI == "" {
		fmt.Fprintln(os.Stderr, "TEST_MONGO_URI not set — skipping POC tests")
		os.Exit(0)
	}

	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "parliamentdb"
	}

	if s := os.Getenv("TEST_LAW_ID"); s != "" {
		if id, err := strconv.Atoi(s); err == nil {
			lawId = id
		}
	}

	database, err := db.Connect(mongoURI, dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect to poc MongoDB: %v\n", err)
		os.Exit(1)
	}

	idxCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.EnsureIndexes(idxCtx, database); err != nil {
		fmt.Fprintf(os.Stderr, "ensure indexes: %v\n", err)
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	partyRepo := repository.NewPartyRepo(database)
	politiciansRepo := repository.NewPoliticiansRepo(database)
	votingRepo := repository.NewVotingRepo(database)

	partyH := handlers.NewPartyHandler(partyRepo)
	politiciansH := handlers.NewPoliticiansHandler(politiciansRepo)
	votingH := handlers.NewVotingHandler(votingRepo)

	router.Get("/health", handlers.Health)
	router.Route("/api/Party", func(r chi.Router) {
		r.Get("/all", partyH.All)
		r.Get("/GetById/", partyH.GetByID)
		r.Get("/query", partyH.Query)
	})
	router.Route("/api/Politicians", func(r chi.Router) {
		r.Get("/getAllPoliticians", politiciansH.GetAll)
		r.Get("/GetById/", politiciansH.GetByID)
		r.Get("/GetByName/", politiciansH.GetByName)
	})
	router.Route("/api/Voting", func(r chi.Router) {
		r.Get("/getAllRounds", votingH.GetAllRounds)
		r.Get("/getRoundById/", votingH.GetRoundById)
		r.Get("/GetResultForVote/", votingH.GetResultForVote)
		r.Get("/GetAllVotesForARoundById/", votingH.GetAllVotesForARoundById)
	})

	pocServer = httptest.NewServer(router)
	defer pocServer.Close()

	os.Exit(m.Run())
}

func TestHealth(t *testing.T) {
	resp, err := http.Get(pocServer.URL + "/health")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestPartiesExist(t *testing.T) {
	resp, err := http.Get(pocServer.URL + "/api/Party/all")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var parties []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&parties); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(parties) == 0 {
		t.Error("expected at least one party after import, got none")
	}
	t.Logf("found %d parties", len(parties))

	for _, party := range parties {
		if party["acronym"] == nil {
			t.Error("party missing 'acronym' field")
		}
		if party["name"] == nil {
			t.Error("party missing 'name' field")
		}
	}
}

func TestPoliticiansExist(t *testing.T) {
	resp, err := http.Get(pocServer.URL + "/api/Politicians/getAllPoliticians")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var politicians []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&politicians); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(politicians) == 0 {
		t.Error("expected at least one politician after import, got none")
	}
	t.Logf("found %d politicians", len(politicians))

	for _, politician := range politicians {
		if politician["name"] == nil {
			t.Error("politician missing 'name' field")
		}
		if politician["partyId"] == nil {
			t.Error("politician missing 'partyId' field")
		}
	}
}

func TestImportedRoundExists(t *testing.T) {
	if lawId == 0 {
		t.Skip("TEST_LAW_ID not set, skipping round lookup")
	}

	url := fmt.Sprintf("%s/api/Voting/getRoundById/?voteNumber=%d", pocServer.URL, lawId)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for lawId %d, got %d — law may not exist on cdep.ro", lawId, resp.StatusCode)
	}

	var round map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&round); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if round["voteId"] == nil {
		t.Error("round missing 'voteId' field")
	}
	if round["title"] == nil {
		t.Error("round missing 'title' field")
	}
	if round["voteDate"] == nil {
		t.Error("round missing 'voteDate' field")
	}

	t.Logf("imported round: voteId=%v title=%v voteDate=%v", round["voteId"], round["title"], round["voteDate"])
}

func TestImportedVotesHaveCorrectStructure(t *testing.T) {
	if lawId == 0 {
		t.Skip("TEST_LAW_ID not set, skipping vote structure check")
	}

	url := fmt.Sprintf("%s/api/Voting/GetResultForVote/?number=%d", pocServer.URL, lawId)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var votes []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&votes); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(votes) == 0 {
		t.Fatalf("expected votes for lawId %d, got none", lawId)
	}
	t.Logf("found %d votes for lawId %d", len(votes), lawId)

	for i, vote := range votes {
		politician, ok := vote["politician"].(map[string]any)
		if !ok {
			t.Errorf("vote[%d] missing nested 'politician' object", i)
			continue
		}
		if politician["name"] == nil {
			t.Errorf("vote[%d] politician missing 'name'", i)
		}
		if _, ok := politician["party"].(map[string]any); !ok {
			t.Errorf("vote[%d] politician missing nested 'party' object", i)
		}
		if vote["round"] == nil {
			t.Errorf("vote[%d] missing 'round' object", i)
		}
		pos, ok := vote["position"].(float64)
		if !ok || pos < 0 || pos > 3 {
			t.Errorf("vote[%d] invalid position: %v", i, vote["position"])
		}
	}
}

func TestAllRoundsEndpointWorks(t *testing.T) {
	resp, err := http.Get(pocServer.URL + "/api/Voting/getAllRounds")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var rounds []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rounds); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(rounds) == 0 {
		t.Error("expected at least one round after import, got none")
	}
	t.Logf("found %d rounds", len(rounds))
}
