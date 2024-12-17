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

func TestListPullRequests(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{
		PullRequests: []*github.PullRequest{
			{Number: github.Int(1)},
			{Number: github.Int(2)},
		},
	}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	req, _ := http.NewRequest("GET", "/repos/test-repo/pulls?n=1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Additional assertions can be added here
}

func Test_ListPullRequests_Success_NoLimit(t *testing.T) {
	mockPulls := []*github.PullRequest{
		{Title: github.String("PR 1")},
		{Title: github.String("PR 2")},
	}
	mockClient := &mocks.MockGitHubClient{PullRequests: mockPulls}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	req, err := http.NewRequest("GET", "/repos/test_repo/pulls", nil)
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

	if response["repository"] != "test_repo" {
		t.Errorf("Expected repository 'test_repo', got '%v'", response["repository"])
	}

	if response["count"] != float64(len(mockPulls)) {
		t.Errorf("Expected count %d, got %v", len(mockPulls), response["count"])
	}

	prs, ok := response["pull_requests"].([]interface{})
	if !ok {
		t.Errorf("Expected 'pull_requests' to be a list, got '%v'", response["pull_requests"])
	}

	if len(prs) != len(mockPulls) {
		t.Errorf("Expected %d pull requests, got %d", len(mockPulls), len(prs))
	}

	for i, pr := range prs {
		prMap, ok := pr.(map[string]interface{})
		if !ok {
			t.Errorf("Expected PR to be a map, got '%v'", pr)
		}

		expectedTitle := *mockPulls[i].Title
		if prMap["title"] != expectedTitle {
			t.Errorf("Expected PR title '%s', got '%v'", expectedTitle, prMap["title"])
		}
	}
}

func Test_ListPullRequests_Success_WithLimit(t *testing.T) {
	mockPulls := []*github.PullRequest{
		{Title: github.String("PR 1")},
		{Title: github.String("PR 2")},
		{Title: github.String("PR 3")},
	}
	mockClient := &mocks.MockGitHubClient{PullRequests: mockPulls}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	req, err := http.NewRequest("GET", "/repos/test_repo/pulls?n=2", nil)
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

	if response["repository"] != "test_repo" {
		t.Errorf("Expected repository 'test_repo', got '%v'", response["repository"])
	}

	count, ok := response["count"].(float64)
	if !ok {
		t.Errorf("Expected 'count' to be a number, got '%v'", response["count"])
	}

	if count != 2 {
		t.Errorf("Expected count 2, got %v", count)
	}

	prs, ok := response["pull_requests"].([]interface{})
	if !ok {
		t.Errorf("Expected 'pull_requests' to be a list, got '%v'", response["pull_requests"])
	}

	if len(prs) != 2 {
		t.Errorf("Expected 2 pull requests, got %d", len(prs))
	}

	expectedTitles := []string{"PR 1", "PR 2"}
	for i, pr := range prs {
		prMap, ok := pr.(map[string]interface{})
		if !ok {
			t.Errorf("Expected PR to be a map, got '%v'", pr)
		}

		expectedTitle := expectedTitles[i]
		if prMap["title"] != expectedTitle {
			t.Errorf("Expected PR title '%s', got '%v'", expectedTitle, prMap["title"])
		}
	}
}

func Test_ListPullRequests_Success_FewerThanLimit(t *testing.T) {
	mockPulls := []*github.PullRequest{
		{Title: github.String("PR 1")},
	}
	mockClient := &mocks.MockGitHubClient{PullRequests: mockPulls}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	req, err := http.NewRequest("GET", "/repos/test_repo/pulls?n=5", nil)
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

	if response["repository"] != "test_repo" {
		t.Errorf("Expected repository 'test_repo', got '%v'", response["repository"])
	}

	count, ok := response["count"].(float64)
	if !ok {
		t.Errorf("Expected 'count' to be a number, got '%v'", response["count"])
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %v", count)
	}

	prs, ok := response["pull_requests"].([]interface{})
	if !ok {
		t.Errorf("Expected 'pull_requests' to be a list, got '%v'", response["pull_requests"])
	}

	if len(prs) != 1 {
		t.Errorf("Expected 1 pull request, got %d", len(prs))
	}

	prMap, ok := prs[0].(map[string]interface{})
	if !ok {
		t.Errorf("Expected PR to be a map, got '%v'", prs[0])
	}

	if prMap["title"] != "PR 1" {
		t.Errorf("Expected PR title 'PR 1', got '%v'", prMap["title"])
	}
}

func Test_ListPullRequests_InvalidN(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	req, err := http.NewRequest("GET", "/repos/test_repo/pulls?n=invalid", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["error"] != "invalid value for n" {
		t.Errorf("Expected error 'invalid value for n', got '%v'", response["error"])
	}
}

func Test_ListPullRequests_ClientError(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{Err: errors.New("mock error")}
	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
	router := SetupRouter(ghClient)

	req, err := http.NewRequest("GET", "/repos/test_repo/pulls?n=2", nil)
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
