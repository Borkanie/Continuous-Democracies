package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestPartyAll_ReturnsActiveParties(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Party/all?active=true")
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

	if len(parties) != 4 {
		t.Errorf("expected 4 active parties, got %d", len(parties))
	}

	for _, party := range parties {
		if party["active"] != true {
			t.Errorf("expected all parties active, got inactive: %v", party["name"])
		}
		if _, ok := party["color"].(string); !ok {
			t.Errorf("expected color to be a string, got %T", party["color"])
		}
	}
}

func TestPartyAll_ReturnsAllWithNoFilter(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Party/all")
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

	if len(parties) != 5 {
		t.Errorf("expected 5 parties, got %d", len(parties))
	}
}

func TestPartyGetById_ExistingParty(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Party/GetById/?id=11111111-1111-1111-1111-111111111001")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var party map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&party); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if party["acronym"] != "PSD" {
		t.Errorf("expected acronym PSD, got %v", party["acronym"])
	}
}

func TestPartyGetById_NotFound(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Party/GetById/?id=99999999-9999-9999-9999-999999999999")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestPartyQuery_ByAcronym(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Party/query?acronym=PNL")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var party map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&party); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if party["acronym"] != "PNL" {
		t.Errorf("expected acronym PNL, got %v", party["acronym"])
	}
}

func TestPartyQuery_MissingParams(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Party/query")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}
