package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetAllPoliticians_ReturnsAll(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Politicians/getAllPoliticians?number=500")
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

	if len(politicians) != 7 {
		t.Errorf("expected 7 politicians, got %d", len(politicians))
	}

	// Verify each politician has a nested party object
	for _, politician := range politicians {
		party, ok := politician["party"].(map[string]any)
		if !ok {
			t.Errorf("politician %v missing nested party object", politician["name"])
			continue
		}
		if party["acronym"] == nil {
			t.Errorf("politician %v party missing acronym", politician["name"])
		}
	}
}

func TestGetAllPoliticians_FilterByPartyAcronym(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Politicians/getAllPoliticians?partyAcronym=PSD&number=500")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var politicians []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&politicians); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	// PSD has Ion Popescu, Maria Ionescu, Andrei Popa (3 total)
	if len(politicians) != 3 {
		t.Errorf("expected 3 PSD politicians, got %d", len(politicians))
	}
}

func TestGetAllPoliticians_FilterByActive(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Politicians/getAllPoliticians?isActive=true&number=500")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var politicians []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&politicians); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	// 6 active politicians (all except Andrei Popa)
	if len(politicians) != 6 {
		t.Errorf("expected 6 active politicians, got %d", len(politicians))
	}
}

func TestGetAllPoliticians_FilterByGender(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Politicians/getAllPoliticians?gender=1&number=500")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var politicians []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&politicians); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	// 3 female politicians: Maria Ionescu, Ana Constantin, Elena Vasile
	if len(politicians) != 3 {
		t.Errorf("expected 3 female politicians, got %d", len(politicians))
	}
}

func TestGetPoliticianById_Exists(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Politicians/GetById/?id=22222222-2222-2222-2222-222222222001")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var politician map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&politician); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if politician["name"] != "Ion Popescu" {
		t.Errorf("expected Ion Popescu, got %v", politician["name"])
	}
}

func TestGetPoliticianById_NotFound(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Politicians/GetById/?id=99999999-9999-9999-9999-999999999999")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestGetPoliticianByName_CaseInsensitive(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Politicians/GetByName/?name=ion+popescu")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var politician map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&politician); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if politician["name"] != "Ion Popescu" {
		t.Errorf("expected Ion Popescu, got %v", politician["name"])
	}
}
