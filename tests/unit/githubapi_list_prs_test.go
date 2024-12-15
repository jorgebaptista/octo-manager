package githubapi_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v67/github"
	"github.com/jorgebaptista/octo-manager/internal/githubapi"
	"github.com/jorgebaptista/octo-manager/tests/mocks"
)

func TestClient_ListPullRequests_Success(t *testing.T) {
	mockPulls := []*github.PullRequest{
		{Title: github.String("PR 1")},
		{Title: github.String("PR 2")},
	}
	mockClient := &mocks.MockGitHubClient{PullRequests: mockPulls}
	client := githubapi.NewTestClient(mockClient, "test_owner")

	prs, err := client.ListPullRequests(context.Background(), "test_repo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(prs) != 2 {
		t.Fatalf("expected 2 pull requests, got %d", len(prs))
	}
	if *prs[0].Title != "PR 1" || *prs[1].Title != "PR 2" {
		t.Errorf("unexpected pull request titles: %+v", prs)
	}
}

func TestClient_ListPullRequests_Error(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{Err: githubapi.Error("mock error")}
	client := githubapi.NewTestClient(mockClient, "test_owner")

	_, err := client.ListPullRequests(context.Background(), "test_repo")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
