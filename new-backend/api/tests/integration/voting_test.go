package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetAllRounds_ReturnsAll(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Voting/getAllRounds")
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

	if len(rounds) != 5 {
		t.Errorf("expected 5 rounds, got %d", len(rounds))
	}

	// Rounds should be sorted by voteDate descending
	if len(rounds) >= 2 {
		// voteDate is returned as string; just check first has higher voteId (fixture order)
		firstVoteId := rounds[0]["voteId"].(float64)
		secondVoteId := rounds[1]["voteId"].(float64)
		if firstVoteId < secondVoteId {
			t.Errorf("rounds should be sorted descending by voteDate, got voteId %v before %v", firstVoteId, secondVoteId)
		}
	}
}

func TestGetAllRounds_KeywordFilter(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Voting/getAllRounds?keywords=buget")
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

	if len(rounds) != 1 {
		t.Errorf("expected 1 round matching 'buget', got %d", len(rounds))
	}
}

func TestGetAllRounds_Limit(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Voting/getAllRounds?maxNumberOfEntries=2")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var rounds []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rounds); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(rounds) != 2 {
		t.Errorf("expected 2 rounds with limit=2, got %d", len(rounds))
	}
}

func TestGetRoundById_Existing(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Voting/getRoundById/?voteNumber=1001")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var round map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&round); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if round["voteId"].(float64) != 1001 {
		t.Errorf("expected voteId 1001, got %v", round["voteId"])
	}

	// Both title and name must be present
	if round["title"] == nil {
		t.Error("missing 'title' field")
	}
	if round["name"] == nil {
		t.Error("missing 'name' field")
	}
}

func TestGetRoundById_NotFound(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Voting/getRoundById/?voteNumber=99999")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestGetRoundById_MissingParam(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Voting/getRoundById/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGetResultForVote_ReturnsVotesWithNestedObjects(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Voting/GetResultForVote/?number=1001")
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

	// 6 votes for round 1001
	if len(votes) != 6 {
		t.Errorf("expected 6 votes for round 1001, got %d", len(votes))
	}

	for _, vote := range votes {
		// Each vote must have politician with nested party
		politician, ok := vote["politician"].(map[string]any)
		if !ok {
			t.Errorf("vote missing 'politician' object")
			continue
		}
		if _, ok := politician["party"].(map[string]any); !ok {
			t.Errorf("politician missing nested 'party' object")
		}
		// Each vote must have round
		if vote["round"] == nil {
			t.Error("vote missing 'round' object")
		}
		// Each vote must have name
		if vote["name"] == nil {
			t.Error("vote missing 'name' field")
		}
		// Position must be 0-3
		pos, ok := vote["position"].(float64)
		if !ok || pos < 0 || pos > 3 {
			t.Errorf("invalid position: %v", vote["position"])
		}
	}
}

func TestGetResultForVote_FilterByPartyId(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Voting/GetResultForVote/?number=1001&partyId=11111111-1111-1111-1111-111111111001")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var votes []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&votes); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	// PSD has Ion Popescu (pos 0), Maria Ionescu (pos 0) in round 1001
	if len(votes) != 2 {
		t.Errorf("expected 2 PSD votes in round 1001, got %d", len(votes))
	}
}

func TestGetResultForVote_UnknownRound(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Voting/GetResultForVote/?number=99999")
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

	if len(votes) != 0 {
		t.Errorf("expected empty array for unknown round, got %d votes", len(votes))
	}
}

func TestGetAllVotesForARoundById(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/Voting/GetAllVotesForARoundById/?roundId=33333333-3333-3333-3333-333333333001")
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

	if len(votes) != 6 {
		t.Errorf("expected 6 votes for round UUID, got %d", len(votes))
	}
}

func TestHealth(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/health")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
