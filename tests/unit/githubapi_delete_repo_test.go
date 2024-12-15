package githubapi_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v67/github"
	"github.com/jorgebaptista/octo-manager/internal/githubapi"
	"github.com/jorgebaptista/octo-manager/tests/mocks"
)

func TestClient_DeleteRepo_Success(t *testing.T) {
	mockRepos := []*github.Repository{{Name: github.String("repo1")}, {Name: github.String("repo2")}}
	mockClient := &mocks.MockGitHubClient{Repos: mockRepos}
	client := githubapi.NewTestClient(mockClient, "test_owner")

	err := client.DeleteRepo(context.Background(), "repo1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(mockClient.Repos) != 1 {
		t.Fatalf("expected 1 repo remaining, got %d", len(mockClient.Repos))
	}
	if *mockClient.Repos[0].Name != "repo2" {
		t.Errorf("expected remaining repo to be %q, got %q", "repo2", *mockClient.Repos[0].Name)
	}
}

func TestClient_DeleteRepo_Error(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{Err: githubapi.Error("mock error")}
	client := githubapi.NewTestClient(mockClient, "test_owner")

	err := client.DeleteRepo(context.Background(), "repo1")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
