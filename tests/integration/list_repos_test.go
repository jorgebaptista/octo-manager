package integration

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v67/github"
	"github.com/jorgebaptista/octo-manager/internal/githubapi"
	"github.com/jorgebaptista/octo-manager/tests/mocks"
)

func Test_ListRepos_Success(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{
		Repos: []*github.Repository{
			{Name: github.String("repo1")},
			{Name: github.String("repo2")},
		},
	}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	req, err := http.NewRequest("GET", "/repos", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	repos, ok := response["repositories"].([]interface{})
	if !ok {
		t.Errorf("Expected 'repositories' to be a list, got '%v'", response["repositories"])
	}

	if len(repos) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(repos))
	}

	if repos[0] != "repo1" {
		t.Errorf("Expected first repository to be 'repo1', got '%v'", repos[0])
	}

	if repos[1] != "repo2" {
		t.Errorf("Expected second repository to be 'repo2', got '%v'", repos[1])
	}
}

func Test_ListRepos_ClientError(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{Err: errors.New("mock error")}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	req, err := http.NewRequest("GET", "/repos", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["error"] != "mock error" {
		t.Errorf("Expected error 'mock error', got '%v'", response["error"])
	}
}
