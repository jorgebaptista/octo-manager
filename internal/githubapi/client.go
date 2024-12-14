package githubapi

import (
	"os"

	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

type Client struct {
	gh    *github.Client
	owner string
}

func NewClient() (*Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	owner := os.Getenv("GITHUB_OWNER")

	if token == "" {
		return nil, ErrMissingToken
	}
	if owner == "" {
		return nil, ErrMissingOwner
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	ghClient := github.NewClient(tc)

	return &Client{
		gh:    ghClient,
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
