package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v67/github"
	"github.com/jorgebaptista/octo-manager/internal/githubapi"
	"github.com/jorgebaptista/octo-manager/tests/mocks"
)

func Test_DeleteRepo_Success(t *testing.T) {
	repoName := "existing-repo"
	mockClient := &mocks.MockGitHubClient{
		Repos: []*github.Repository{
			{Name: github.String(repoName)},
		},
	}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	// Simulate DELETE /repos/existing-repo
	req, err := http.NewRequest("DELETE", "/repos/"+repoName, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Assert response content
	if response["message"] != "Repository deleted" {
		t.Errorf("Expected message 'Repository deleted', got '%v'", response["message"])
	}

	if response["repo"] != repoName {
		t.Errorf("Expected repo name '%s', got '%v'", repoName, response["repo"])
	}

	// Assert that the repository was removed from the mock client
	if len(mockClient.Repos) != 0 {
		t.Errorf("Expected 0 repositories, got %d", len(mockClient.Repos))
	}
}

func Test_DeleteRepo_NotFound(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{} // No repositories
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	// Simulate DELETE /repos/non-existent-repo
	req, err := http.NewRequest("DELETE", "/repos/non-existent-repo", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

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
	if response["error"] != "repository not found" {
		t.Errorf("Expected error 'repository not found', got '%v'", response["error"])
	}

	// Assert that no repository was removed from the mock client
	if len(mockClient.Repos) != 0 {
		t.Errorf("Expected 0 repositories, got %d", len(mockClient.Repos))
	}
}
