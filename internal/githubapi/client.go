package githubapi

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

type GitHubClient interface {
	ListReposForOwner(ctx context.Context, owner string) ([]*github.Repository, error)
	CreateRepoForOwner(ctx context.Context, owner, repoName string) (*github.Repository, error)
	DeleteRepoForOwner(ctx context.Context, owner, repoName string) error
	ListPullRequestsForOwner(ctx context.Context, owner, repoName string, n int) ([]*github.PullRequest, error)
}

// Real implementation of the GitHubClient interface
type RealGitHubClient struct {
	gh *github.Client
}

// todo log errors?
func (r *RealGitHubClient) CreateRepoForOwner(ctx context.Context, owner, repoName string) (*github.Repository, error) {
	newRepo := &github.Repository{Name: github.String(repoName)}
	repo, _, err := r.gh.Repositories.Create(ctx, owner, newRepo)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *RealGitHubClient) DeleteRepoForOwner(ctx context.Context, owner, repoName string) error {
	repos, _, err := r.gh.Repositories.ListByAuthenticatedUser(ctx, nil)
	if err != nil {
		return err
	}

	// Check if the repo exists
	repoFound := false
	for _, repo := range repos {
		if *repo.Name == repoName {
			repoFound = true
			break
		}
	}

	if !repoFound {
		return fmt.Errorf("repository not found")
	}

	_, err = r.gh.Repositories.Delete(ctx, owner, repoName)
	if err != nil {
		return err
	}

	return nil
}

func (r *RealGitHubClient) ListReposForOwner(ctx context.Context, owner string) ([]*github.Repository, error) {
	repos, _, err := r.gh.Repositories.ListByAuthenticatedUser(ctx, nil)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (r *RealGitHubClient) ListPullRequestsForOwner(ctx context.Context, owner, repoName string, n int) ([]*github.PullRequest, error) {
	opts := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allPRs []*github.PullRequest
	for {
		prs, resp, err := r.gh.PullRequests.List(ctx, owner, repoName, opts)
		if err != nil {
			return nil, err
		}

		allPRs = append(allPRs, prs...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Limit the number of PRs if n is specified
	if n != -1 && len(allPRs) > n {
		allPRs = allPRs[:n]
	}

	return allPRs, nil
}

type Client struct {
	gh    GitHubClient
	owner string
}

func NewClient() (*Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	owner := os.Getenv("GITHUB_OWNER")

	if token == "" {
		log.Println(ErrMissingToken)
		return nil, ErrMissingToken
	}
	if owner == "" {
		log.Println(ErrMissingOwner)
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

func (c *Client) CreateRepo(ctx context.Context, repoName string) (*github.Repository, error) {
	return c.gh.CreateRepoForOwner(ctx, c.owner, repoName)
}

func (c *Client) DeleteRepo(ctx context.Context, repoName string) error {
	return c.gh.DeleteRepoForOwner(ctx, c.owner, repoName)
}

func (c *Client) ListRepos(ctx context.Context) ([]*github.Repository, error) {
	return c.gh.ListReposForOwner(ctx, c.owner)
}

func (c *Client) ListPullRequests(ctx context.Context, repoName string, n int) ([]*github.PullRequest, error) {
	return c.gh.ListPullRequestsForOwner(ctx, c.owner, repoName, n)
}

func NewTestClient(mockClient GitHubClient, owner string) *Client {
	return &Client{
		gh:    mockClient,
		owner: owner,
	}
}
