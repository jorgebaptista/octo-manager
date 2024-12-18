package mocks

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
)

type MockGitHubClient struct {
	Repos        []*github.Repository  // Mock repos
	PullRequests []*github.PullRequest // Mock pull requests
	Err          error
}

func (m *MockGitHubClient) CreateRepoForOwner(ctx context.Context, owner, repoName string) (*github.Repository, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	repo := &github.Repository{Name: github.String(repoName)}
	m.Repos = append(m.Repos, repo) // Add to mock repos
	return repo, nil
}

func (m *MockGitHubClient) DeleteRepoForOwner(ctx context.Context, owner, repoName string) error {
	if m.Err != nil {
		return m.Err
	}

	// Check if the repo exists
	repoFound := false
	for _, repo := range m.Repos {
		if *repo.Name == repoName {
			repoFound = true
			break
		}
	}

	if !repoFound {
		return fmt.Errorf("repository not found")
	}

	for i, repo := range m.Repos {
		if *repo.Name == repoName {
			m.Repos = append(m.Repos[:i], m.Repos[i+1:]...)
			break
		}
	}

	return nil
}

func (m *MockGitHubClient) ListReposForOwner(ctx context.Context, owner string) ([]*github.Repository, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Repos, nil
}

func (m *MockGitHubClient) ListPullRequestsForOwner(ctx context.Context, owner, repoName string, n int) ([]*github.PullRequest, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if n == -1 {
		return m.PullRequests, nil
	}

	if len(m.PullRequests) > n {
		return m.PullRequests[:n], nil
	}

	return m.PullRequests, nil
}
