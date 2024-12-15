package githubapi_test

import (
	"os"
	"testing"

	"github.com/jorgebaptista/octo-manager/internal/githubapi"
)

func TestNewClient_MissingToken(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "") // empty
	os.Setenv("GITHUB_OWNER", "test_owner")

	// Attempt to create a client with missing token
	c, err := githubapi.NewClient()
	if err == nil {
		// If there was no error
		t.Fatal("expected an error when token is missing, got nil")
	}
	if c != nil {
		// If client isnt empty
		t.Fatal("expected no client when token is missing")
	}
}

func TestNewClient_MissingOwner(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "test_token")
	os.Setenv("GITHUB_OWNER", "") // empty

	c, err := githubapi.NewClient()
	if err == nil {
		t.Fatal("expected an error when owner is missing, got nil")
	}
	if c != nil {
		t.Fatal("expected no client when owner is missing")
	}
}

func TestNewClient_Success(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "test_token")
	os.Setenv("GITHUB_OWNER", "test_owner")

	c, err := githubapi.NewClient()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c == nil {
		t.Fatal("expected a client, got nil")
	}
}
