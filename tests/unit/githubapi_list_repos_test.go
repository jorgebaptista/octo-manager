package githubapi_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v67/github"
	"github.com/jorgebaptista/octo-manager/internal/githubapi"
	"github.com/jorgebaptista/octo-manager/tests/mocks"
)

func TestClient_ListRepos_Success(t *testing.T) {
	// Mock data
	mockRepos := []*github.Repository{
		{Name: github.String("repo1")},
		{Name: github.String("repo2")},
	}

	mockClient := &mocks.MockGitHubClient{Repos: mockRepos}
	c := githubapi.NewTestClient(mockClient, "test_owner")

	repos, err := c.ListRepos(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repos) != 2 {
		t.Fatalf("expected 2 repos, got %d", len(repos))
	}
	if *repos[0].Name != "repo1" || *repos[1].Name != "repo2" {
		t.Errorf("unexpected repo names: got %+v", repos)
	}
}

func TestClient_ListRepos_Error(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{
		Err: githubapi.Error("test error"),
	}
	c := githubapi.NewTestClient(mockClient, "test_owner")

	repos, err := c.ListRepos(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if repos != nil {
		t.Fatalf("expected no repos on error, got %+v", repos)
	}
}
