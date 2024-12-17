package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jorgebaptista/octo-manager/internal/githubapi"
	"github.com/jorgebaptista/octo-manager/tests/mocks"
)

func Test_CreateRepo_Success(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	// Simulate POST /repos with valid JSON body
	repoName := "new-repo"
	reqBody := `{"name": "` + repoName + `"}`

	req, err := http.NewRequest("POST", "/repos", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Assert response content
	if response["message"] != "Repository created" {
		t.Errorf("Expected message 'Repository created', got '%v'", response["message"])
	}

	if response["name"] != repoName {
		t.Errorf("Expected repository name '%s', got '%v'", repoName, response["name"])
	}

	if len(mockClient.Repos) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(mockClient.Repos))
	}

	if *mockClient.Repos[0].Name != repoName {
		t.Errorf("Expected repository name '%s', got '%s'", repoName, *mockClient.Repos[0].Name)
	}
}

func Test_CreateRepo_InvalidRequest(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	// Simulate POST /repos with invalid JSON body
	reqBody := `{"invalid": "data"}`

	req, err := http.NewRequest("POST", "/repos", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Assert response content
	if response["error"] != "invalid request" {
		t.Errorf("Expected error 'invalid request', got '%v'", response["error"])
	}

	// Assert that no repository was added to the mock client
	if len(mockClient.Repos) != 0 {
		t.Errorf("Expected 0 repositories, got %d", len(mockClient.Repos))
	}
}

func Test_CreateRepo_ClientError(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{Err: errors.New("mock error")}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	// Simulate POST /repos with valid JSON body
	repoName := "new-repo"
	reqBody := `{"name": "` + repoName + `"}`

	req, err := http.NewRequest("POST", "/repos", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response status code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Assert response content
	if response["error"] != "mock error" {
		t.Errorf("Expected error 'mock error', got '%v'", response["error"])
	}

	// Assert that no repository was added to the mock client
	if len(mockClient.Repos) != 0 {
		t.Errorf("Expected 0 repositories, got %d", len(mockClient.Repos))
	}
}
