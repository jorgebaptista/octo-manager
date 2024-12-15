package githubapi

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

type GitHubClient interface {
	ListReposForOwner(ctx context.Context, owner string) ([]*github.Repository, error)
}

// Real implementation of the GitHubClient interface
type RealGitHubClient struct {
	gh *github.Client
}

func (r *RealGitHubClient) ListReposForOwner(ctx context.Context, owner string) ([]*github.Repository, error) {
	repos, _, err := r.gh.Repositories.ListByAuthenticatedUser(ctx, nil)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

type Client struct {
	gh    GitHubClient
	owner string
}

func NewClient() (*Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	owner := os.Getenv("GITHUB_OWNER")

	if token == "" {
		log.Println("GITHUB_TOKEN is not set")
		return nil, ErrMissingToken
	}
	if owner == "" {
		log.Println("GITHUB_OWNER is not set")
		return nil, ErrMissingOwner
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	ghClient := github.NewClient(tc)
	real := &RealGitHubClient{gh: ghClient}

	return &Client{
		gh:    real,
		owner: owner,
	}, nil

}

// Define custom error types
var (
	ErrMissingToken = Error("GITHUB_TOKEN not set")
	ErrMissingOwner = Error("GITHUB_OWNER not set")
)

// Create a custom error type to implement the error interface
type Error string

func (e Error) Error() string { return string(e) }

func (c *Client) ListRepos(ctx context.Context) ([]*github.Repository, error) {
	return c.gh.ListReposForOwner(ctx, c.owner)
}

func NewTestClient(mockClient GitHubClient, owner string) *Client {
	return &Client{
		gh:    mockClient,
		owner: owner,
	}
}
