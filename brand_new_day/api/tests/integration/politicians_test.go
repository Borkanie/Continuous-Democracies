package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListPoliticians_ReturnsAllPoliticians(testingInstance *testing.T) {
	listPoliticiansResponse, err := http.Get(testServer.URL + "/politicians")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer listPoliticiansResponse.Body.Close()

	if listPoliticiansResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", listPoliticiansResponse.StatusCode)
	}

	var politiciansList []map[string]interface{}
	if err := json.NewDecoder(listPoliticiansResponse.Body).Decode(&politiciansList); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	if len(politiciansList) != 3 {
		testingInstance.Errorf("expected 3 politicians, got %d", len(politiciansList))
	}

	for _, politicianData := range politiciansList {
		if _, hasID := politicianData["id"]; !hasID {
			testingInstance.Error("expected id field in politician")
		}
		if _, hasName := politicianData["name"]; !hasName {
			testingInstance.Error("expected name field in politician")
		}
		if _, hasPartyID := politicianData["partyId"]; !hasPartyID {
			testingInstance.Error("expected partyId field in politician")
		}
		if _, hasActive := politicianData["active"]; !hasActive {
			testingInstance.Error("expected active field in politician")
		}
	}
}

func TestGetPoliticianByID_ExistingPolitician(testingInstance *testing.T) {
	getPoliticianResponse, err := http.Get(testServer.URL + "/politicians/politician-001")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getPoliticianResponse.Body.Close()

	if getPoliticianResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", getPoliticianResponse.StatusCode)
	}

	var politicianData map[string]interface{}
	if err := json.NewDecoder(getPoliticianResponse.Body).Decode(&politicianData); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	if politicianData["name"] != "John Smith" {
		testingInstance.Errorf("expected name 'John Smith', got %v", politicianData["name"])
	}

	if politicianData["partyId"] != "party-001" {
		testingInstance.Errorf("expected partyId 'party-001', got %v", politicianData["partyId"])
	}
}

func TestGetPoliticianByID_NotFound(testingInstance *testing.T) {
	getPoliticianResponse, err := http.Get(testServer.URL + "/politicians/nonexistent-politician")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getPoliticianResponse.Body.Close()

	if getPoliticianResponse.StatusCode != http.StatusNotFound {
		testingInstance.Errorf("expected 404, got %d", getPoliticianResponse.StatusCode)
	}

	var errorResponse map[string]interface{}
	if err := json.NewDecoder(getPoliticianResponse.Body).Decode(&errorResponse); err != nil {
		testingInstance.Fatalf("decode error response: %v", err)
	}

	if _, hasError := errorResponse["error"]; !hasError {
		testingInstance.Error("expected error field in response")
	}
}

func TestGetVotesByPolitician_ReturnsVotes(testingInstance *testing.T) {
	getVotesResponse, err := http.Get(testServer.URL + "/politicians/politician-001/votes")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getVotesResponse.Body.Close()

	if getVotesResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", getVotesResponse.StatusCode)
	}

	var voteEntries []map[string]interface{}
	if err := json.NewDecoder(getVotesResponse.Body).Decode(&voteEntries); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	// politician-001 voted in rounds 1, 2, 3, and 4
	if len(voteEntries) != 4 {
		testingInstance.Errorf("expected 4 vote entries, got %d", len(voteEntries))
	}

	// Verify that votes are sorted by voteDate descending
	// Round 3 (2023-02-10), Round 2 (2023-02-05), Round 1 (2023-02-01)
	if len(voteEntries) >= 2 {
		firstVoteDate := voteEntries[0]["votingRound"].(map[string]interface{})["voteDate"]
		secondVoteDate := voteEntries[1]["votingRound"].(map[string]interface{})["voteDate"]

		if firstVoteDate.(string) < secondVoteDate.(string) {
			testingInstance.Error("expected votes to be sorted by voteDate descending")
		}
	}

	// Check that votingRound.votes is empty in politician vote entries
	for _, voteEntry := range voteEntries {
		votingRound := voteEntry["votingRound"].(map[string]interface{})
		votes := votingRound["votes"].([]interface{})
		if len(votes) != 0 {
			testingInstance.Error("expected votingRound.votes to be empty in politician vote entries")
		}
	}

	// Check that normative and lawBucket are present
	for _, voteEntry := range voteEntries {
		if _, hasValue := voteEntry["value"]; !hasValue {
			testingInstance.Error("expected value field in vote entry")
		}
		if _, hasVotingRound := voteEntry["votingRound"]; !hasVotingRound {
			testingInstance.Error("expected votingRound field in vote entry")
		}
		if _, hasNormative := voteEntry["normative"]; !hasNormative {
			testingInstance.Error("expected normative field in vote entry")
		}
		if _, hasLawBucket := voteEntry["lawBucket"]; !hasLawBucket {
			testingInstance.Error("expected lawBucket field in vote entry")
		}
	}
}

func TestGetVotesByPolitician_NotFound(testingInstance *testing.T) {
	getVotesResponse, err := http.Get(testServer.URL + "/politicians/nonexistent-politician/votes")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getVotesResponse.Body.Close()

	if getVotesResponse.StatusCode != http.StatusNotFound {
		testingInstance.Errorf("expected 404, got %d", getVotesResponse.StatusCode)
	}
}
