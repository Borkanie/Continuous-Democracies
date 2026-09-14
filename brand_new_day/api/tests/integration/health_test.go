package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetHealth_ReturnsOkStatus(testingInstance *testing.T) {
	healthResponse, err := http.Get(testServer.URL + "/health")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer healthResponse.Body.Close()

	if healthResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", healthResponse.StatusCode)
	}

	var healthStatusResponse map[string]string
	if err := json.NewDecoder(healthResponse.Body).Decode(&healthStatusResponse); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	if healthStatusResponse["status"] != "ok" {
		testingInstance.Errorf("expected status 'ok', got %v", healthStatusResponse["status"])
	}
}
