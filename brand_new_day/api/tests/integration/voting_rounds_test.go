package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListVotingRounds_ReturnsAllRounds(testingInstance *testing.T) {
	listRoundsResponse, err := http.Get(testServer.URL + "/votingRounds")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer listRoundsResponse.Body.Close()

	if listRoundsResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", listRoundsResponse.StatusCode)
	}

	var roundsList []map[string]interface{}
	if err := json.NewDecoder(listRoundsResponse.Body).Decode(&roundsList); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	if len(roundsList) != 4 {
		testingInstance.Errorf("expected 4 voting rounds, got %d", len(roundsList))
	}

	// Check that votes arrays are empty in list response
	for _, roundData := range roundsList {
		votesInterface := roundData["votes"]
		if votesInterface == nil {
			testingInstance.Error("expected votes field to be non-nil")
			continue
		}
		votes := votesInterface.([]interface{})
		if len(votes) != 0 {
			testingInstance.Error("expected empty votes array in list response")
		}
	}
}

func TestGetVotingRoundByID_ReturnsDetailWithVotes(testingInstance *testing.T) {
	getRoundResponse, err := http.Get(testServer.URL + "/votingRounds/1")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getRoundResponse.Body.Close()

	if getRoundResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", getRoundResponse.StatusCode)
	}

	var roundData map[string]interface{}
	if err := json.NewDecoder(getRoundResponse.Body).Decode(&roundData); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	if roundData["title"] != "Final vote on budget" {
		testingInstance.Errorf("expected title 'Final vote on budget', got %v", roundData["title"])
	}

	// Verify that the detail response includes the full votes array
	votesInterface := roundData["votes"]
	if votesInterface == nil {
		testingInstance.Fatal("expected votes field to be non-nil")
	}

	votes := votesInterface.([]interface{})
	if len(votes) != 3 {
		testingInstance.Errorf("expected 3 votes, got %d", len(votes))
	}

	// Verify that votes have the correct structure
	for _, voteInterface := range votes {
		voteData := voteInterface.(map[string]interface{})
		if _, hasPoliticianID := voteData["politicianId"]; !hasPoliticianID {
			testingInstance.Error("expected politicianId in vote")
		}
		if _, hasPartyID := voteData["partyId"]; !hasPartyID {
			testingInstance.Error("expected partyId in vote")
		}
		if _, hasValue := voteData["value"]; !hasValue {
			testingInstance.Error("expected value in vote")
		}
	}

	// Verify that normative and lawBucket are embedded
	if _, hasNormative := roundData["normative"]; !hasNormative {
		testingInstance.Error("expected normative field in round detail")
	}

	if _, hasLawBucket := roundData["lawBucket"]; !hasLawBucket {
		testingInstance.Error("expected lawBucket field in round detail")
	}
}

func TestGetVotingRoundByID_ResolveExactNormativeVersion(testingInstance *testing.T) {
	// This is the critical test: round 1 voted on norm-001 version 1,
	// round 3 voted on norm-001 version 2. We verify that each round
	// resolves the exact version it voted on, not the latest version.

	// Get round 1 (voted on norm-001 version 1)
	getRound1Response, err := http.Get(testServer.URL + "/votingRounds/1")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getRound1Response.Body.Close()

	var round1Data map[string]interface{}
	if err := json.NewDecoder(getRound1Response.Body).Decode(&round1Data); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	normative1 := round1Data["normative"].(map[string]interface{})
	if normative1["version"].(float64) != 1 {
		testingInstance.Errorf("expected round 1 to resolve normative version 1, got %v", normative1["version"])
	}
	if normative1["text"] != "The state ensures the implementation of the budget." {
		testingInstance.Errorf("expected original text for version 1, got %v", normative1["text"])
	}

	// Get round 3 (voted on norm-001 version 2)
	getRound3Response, err := http.Get(testServer.URL + "/votingRounds/3")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getRound3Response.Body.Close()

	var round3Data map[string]interface{}
	if err := json.NewDecoder(getRound3Response.Body).Decode(&round3Data); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	normative3 := round3Data["normative"].(map[string]interface{})
	if normative3["version"].(float64) != 2 {
		testingInstance.Errorf("expected round 3 to resolve normative version 2, got %v", normative3["version"])
	}
	if normative3["text"] != "The state ensures the implementation of the budget with oversight." {
		testingInstance.Errorf("expected updated text for version 2, got %v", normative3["text"])
	}

	// Now verify that GET /lawBuckets/100 still returns version 2 (the latest)
	getLawBucketResponse, err := http.Get(testServer.URL + "/lawBuckets/100")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getLawBucketResponse.Body.Close()

	var lawBucketData map[string]interface{}
	if err := json.NewDecoder(getLawBucketResponse.Body).Decode(&lawBucketData); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	normativesInterface := lawBucketData["normatives"].([]interface{})
	norm001InBucket := findNormativeByID(normativesInterface, "norm-001")
	if norm001InBucket == nil {
		testingInstance.Fatal("expected to find norm-001 in law bucket")
	}

	if norm001InBucket["version"].(float64) != 2 {
		testingInstance.Errorf("expected law bucket to embed normative version 2, got %v", norm001InBucket["version"])
	}
	if norm001InBucket["text"] != "The state ensures the implementation of the budget with oversight." {
		testingInstance.Errorf("expected law bucket to have updated text, got %v", norm001InBucket["text"])
	}
}

func TestGetVotesByVotingRound_ReturnsHydratedVotes(testingInstance *testing.T) {
	getVotesResponse, err := http.Get(testServer.URL + "/votingRounds/1/votes")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getVotesResponse.Body.Close()

	if getVotesResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", getVotesResponse.StatusCode)
	}

	var hydratedVotesList []map[string]interface{}
	if err := json.NewDecoder(getVotesResponse.Body).Decode(&hydratedVotesList); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	if len(hydratedVotesList) != 3 {
		testingInstance.Errorf("expected 3 hydrated votes, got %d", len(hydratedVotesList))
	}

	// Verify that each vote has politician and party data
	for _, hydratedVoteData := range hydratedVotesList {
		if _, hasPoliticianID := hydratedVoteData["politicianId"]; !hasPoliticianID {
			testingInstance.Error("expected politicianId in hydrated vote")
		}
		if _, hasPartyID := hydratedVoteData["partyId"]; !hasPartyID {
			testingInstance.Error("expected partyId in hydrated vote")
		}
		if _, hasValue := hydratedVoteData["value"]; !hasValue {
			testingInstance.Error("expected value in hydrated vote")
		}

		// Verify politician object is populated
		politicianInterface := hydratedVoteData["politician"]
		if politicianInterface == nil {
			testingInstance.Error("expected politician object to be non-nil")
		} else {
			politicianData := politicianInterface.(map[string]interface{})
			if _, hasName := politicianData["name"]; !hasName {
				testingInstance.Error("expected name in politician object")
			}
		}

		// Verify party object is populated
		partyInterface := hydratedVoteData["party"]
		if partyInterface == nil {
			testingInstance.Error("expected party object to be non-nil")
		} else {
			partyData := partyInterface.(map[string]interface{})
			if _, hasAcronym := partyData["acronym"]; !hasAcronym {
				testingInstance.Error("expected acronym in party object")
			}
		}
	}

	// Verify that the votes maintain order from embedded array
	firstVote := hydratedVotesList[0]
	if firstVote["politicianId"] != "politician-001" {
		testingInstance.Errorf("expected first vote from politician-001, got %v", firstVote["politicianId"])
	}
}

func TestGetVotesByVotingRound_NotFound(testingInstance *testing.T) {
	getVotesResponse, err := http.Get(testServer.URL + "/votingRounds/99999/votes")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getVotesResponse.Body.Close()

	if getVotesResponse.StatusCode != http.StatusNotFound {
		testingInstance.Errorf("expected 404, got %d", getVotesResponse.StatusCode)
	}
}

func TestGetVotingRoundByID_NotFound(testingInstance *testing.T) {
	getRoundResponse, err := http.Get(testServer.URL + "/votingRounds/99999")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getRoundResponse.Body.Close()

	if getRoundResponse.StatusCode != http.StatusNotFound {
		testingInstance.Errorf("expected 404, got %d", getRoundResponse.StatusCode)
	}

	var errorResponse map[string]interface{}
	if err := json.NewDecoder(getRoundResponse.Body).Decode(&errorResponse); err != nil {
		testingInstance.Fatalf("decode error response: %v", err)
	}

	if _, hasError := errorResponse["error"]; !hasError {
		testingInstance.Error("expected error field in response")
	}
}

// Helper function to find a normative by ID in a list of normatives
func findNormativeByID(normativesList []interface{}, normativeID string) map[string]interface{} {
	for _, normativeInterface := range normativesList {
		normativeData := normativeInterface.(map[string]interface{})
		if normativeData["id"].(string) == normativeID {
			return normativeData
		}
	}
	return nil
}
