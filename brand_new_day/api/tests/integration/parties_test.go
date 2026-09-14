package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListParties_ReturnsAllParties(testingInstance *testing.T) {
	listPartiesResponse, err := http.Get(testServer.URL + "/parties")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer listPartiesResponse.Body.Close()

	if listPartiesResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", listPartiesResponse.StatusCode)
	}

	var partiesList []map[string]interface{}
	if err := json.NewDecoder(listPartiesResponse.Body).Decode(&partiesList); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	if len(partiesList) != 3 {
		testingInstance.Errorf("expected 3 parties, got %d", len(partiesList))
	}

	for _, partyData := range partiesList {
		if _, hasID := partyData["id"]; !hasID {
			testingInstance.Error("expected id field in party")
		}
		if _, hasName := partyData["name"]; !hasName {
			testingInstance.Error("expected name field in party")
		}
		if _, hasAcronym := partyData["acronym"]; !hasAcronym {
			testingInstance.Error("expected acronym field in party")
		}
		if _, hasActive := partyData["active"]; !hasActive {
			testingInstance.Error("expected active field in party")
		}
	}
}

func TestGetPartyByID_ExistingParty(testingInstance *testing.T) {
	getPartyResponse, err := http.Get(testServer.URL + "/parties/party-001")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getPartyResponse.Body.Close()

	if getPartyResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", getPartyResponse.StatusCode)
	}

	var partyData map[string]interface{}
	if err := json.NewDecoder(getPartyResponse.Body).Decode(&partyData); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	if partyData["acronym"] != "PSD" {
		testingInstance.Errorf("expected acronym 'PSD', got %v", partyData["acronym"])
	}

	if partyData["name"] != "Social Democratic Party" {
		testingInstance.Errorf("expected name 'Social Democratic Party', got %v", partyData["name"])
	}
}

func TestGetPartyByID_NotFound(testingInstance *testing.T) {
	getPartyResponse, err := http.Get(testServer.URL + "/parties/nonexistent-party")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getPartyResponse.Body.Close()

	if getPartyResponse.StatusCode != http.StatusNotFound {
		testingInstance.Errorf("expected 404, got %d", getPartyResponse.StatusCode)
	}

	var errorResponse map[string]interface{}
	if err := json.NewDecoder(getPartyResponse.Body).Decode(&errorResponse); err != nil {
		testingInstance.Fatalf("decode error response: %v", err)
	}

	if _, hasError := errorResponse["error"]; !hasError {
		testingInstance.Error("expected error field in response")
	}
}
