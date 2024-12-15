package githubapi_test

import (
	"context"
	"testing"

	"github.com/jorgebaptista/octo-manager/internal/githubapi"
	"github.com/jorgebaptista/octo-manager/tests/mocks"
)

func TestClient_CreateRepo_Success(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{}
	client := githubapi.NewTestClient(mockClient, "test_owner")

	repoName := "new-repo"
	repo, err := client.CreateRepo(context.Background(), repoName)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if *repo.Name != repoName {
		t.Errorf("expected repo name %q, got %q", repoName, *repo.Name)
	}
}

func TestClient_CreateRepo_Error(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{
		Err: githubapi.Error("mock error"),
	}
	client := githubapi.NewTestClient(mockClient, "test_owner")

	_, err := client.CreateRepo(context.Background(), "new-repo")

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
