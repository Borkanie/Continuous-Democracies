package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListLawBuckets_ReturnsAllBuckets(testingInstance *testing.T) {
	listBucketsResponse, err := http.Get(testServer.URL + "/lawBuckets")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer listBucketsResponse.Body.Close()

	if listBucketsResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", listBucketsResponse.StatusCode)
	}

	var bucketsList []map[string]interface{}
	if err := json.NewDecoder(listBucketsResponse.Body).Decode(&bucketsList); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	if len(bucketsList) != 2 {
		testingInstance.Errorf("expected 2 law buckets, got %d", len(bucketsList))
	}

	for _, bucketData := range bucketsList {
		if _, hasID := bucketData["id"]; !hasID {
			testingInstance.Error("expected id field in law bucket")
		}
		if _, hasTitle := bucketData["title"]; !hasTitle {
			testingInstance.Error("expected title field in law bucket")
		}
		if _, hasNormatives := bucketData["normatives"]; !hasNormatives {
			testingInstance.Error("expected normatives field in law bucket")
		}
	}
}

func TestGetLawBucketByID_ExistingBucket(testingInstance *testing.T) {
	getBucketResponse, err := http.Get(testServer.URL + "/lawBuckets/100")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getBucketResponse.Body.Close()

	if getBucketResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", getBucketResponse.StatusCode)
	}

	var bucketData map[string]interface{}
	if err := json.NewDecoder(getBucketResponse.Body).Decode(&bucketData); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	if bucketData["title"] != "Law on Budget" {
		testingInstance.Errorf("expected title 'Law on Budget', got %v", bucketData["title"])
	}

	if bucketData["plNumber"] != "PL-001" {
		testingInstance.Errorf("expected plNumber 'PL-001', got %v", bucketData["plNumber"])
	}

	// Verify that normatives are embedded
	normativesInterface := bucketData["normatives"]
	if normativesInterface == nil {
		testingInstance.Fatal("expected normatives field to be non-nil")
	}

	normatives := normativesInterface.([]interface{})
	if len(normatives) == 0 {
		testingInstance.Fatal("expected at least one normative in law bucket")
	}

	// Verify that we get the latest version of each normative
	// norm-001 should have version 2, norm-002 should have version 1
	normativesMap := make(map[string]float64)
	for _, normativeInterface := range normatives {
		normativeData := normativeInterface.(map[string]interface{})
		normativeID := normativeData["id"].(string)
		normativeVersion := normativeData["version"].(float64)
		normativesMap[normativeID] = normativeVersion
	}

	if normativesMap["norm-001"] != 2 {
		testingInstance.Errorf("expected norm-001 version 2, got %v", normativesMap["norm-001"])
	}

	if normativesMap["norm-002"] != 1 {
		testingInstance.Errorf("expected norm-002 version 1, got %v", normativesMap["norm-002"])
	}
}

func TestGetLawBucketByID_NotFound(testingInstance *testing.T) {
	getBucketResponse, err := http.Get(testServer.URL + "/lawBuckets/99999")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getBucketResponse.Body.Close()

	if getBucketResponse.StatusCode != http.StatusNotFound {
		testingInstance.Errorf("expected 404, got %d", getBucketResponse.StatusCode)
	}

	var errorResponse map[string]interface{}
	if err := json.NewDecoder(getBucketResponse.Body).Decode(&errorResponse); err != nil {
		testingInstance.Fatalf("decode error response: %v", err)
	}

	if _, hasError := errorResponse["error"]; !hasError {
		testingInstance.Error("expected error field in response")
	}
}

func TestGetVotingRoundsByLawBucket_ReturnsAllRounds(testingInstance *testing.T) {
	getRoundsResponse, err := http.Get(testServer.URL + "/lawBuckets/100/votingRounds")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getRoundsResponse.Body.Close()

	if getRoundsResponse.StatusCode != http.StatusOK {
		testingInstance.Fatalf("expected 200, got %d", getRoundsResponse.StatusCode)
	}

	var roundsList []map[string]interface{}
	if err := json.NewDecoder(getRoundsResponse.Body).Decode(&roundsList); err != nil {
		testingInstance.Fatalf("decode response: %v", err)
	}

	// Law bucket 100 has 3 voting rounds (1, 2, 3)
	if len(roundsList) != 3 {
		testingInstance.Errorf("expected 3 voting rounds for law bucket 100, got %d", len(roundsList))
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
			testingInstance.Errorf("expected empty votes array in list response, got %d votes", len(votes))
		}
	}
}

func TestGetVotingRoundsByLawBucket_NotFound(testingInstance *testing.T) {
	getRoundsResponse, err := http.Get(testServer.URL + "/lawBuckets/99999/votingRounds")
	if err != nil {
		testingInstance.Fatalf("request failed: %v", err)
	}
	defer getRoundsResponse.Body.Close()

	if getRoundsResponse.StatusCode != http.StatusNotFound {
		testingInstance.Errorf("expected 404, got %d", getRoundsResponse.StatusCode)
	}
}
