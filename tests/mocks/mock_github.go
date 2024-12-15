package mocks

import (
	"context"

	"github.com/google/go-github/v67/github"
)

type MockGitHubClient struct {
	Repos []*github.Repository // Mock repos
	Err   error
}

func (m *MockGitHubClient) ListReposForOwner(ctx context.Context, owner string) ([]*github.Repository, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Repos, nil
}
