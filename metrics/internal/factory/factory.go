package factory

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Factory is used to create dependencies for the CLI application
type Factory struct {
	NewEncoder             func() (export.Encoder, error)
	NewExporter            func() (export.Exporter, error)
	NewGitHubHTTPClient    func() (*http.Client, error)
	NewGitHubGraphQLClient func() (github.GraphQLClient, error)
	NewGitHubRepo          func() (*github.Repo, error)
}

// NewFactory creates the default Factory for the CLI application
func NewFactory(ctx context.Context) *Factory {
	f := new(Factory)

	f.NewEncoder = func() (export.Encoder, error) {
		return &export.JSONEcoder{}, nil
	}

	f.NewExporter = func() (export.Exporter, error) {
		encoder, err := f.NewEncoder()
		if err != nil {
			return nil, err
		}
		return &export.WriterExporter{W: os.Stdout, Encoder: encoder}, nil
	}

	f.NewGitHubHTTPClient = func() (*http.Client, error) {
		token := config.FromEnv("GITHUB", "TOKEN")
		if token == "" {
			return nil, fmt.Errorf("Requires GitHub token to be set in env")
		}
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		return oauth2.NewClient(ctx, src), nil
	}

	f.NewGitHubGraphQLClient = func() (github.GraphQLClient, error) {
		httpClient, err := f.NewGitHubHTTPClient()
		if err != nil {
			return nil, err
		}
		gqlClient := githubv4.NewClient(httpClient)
		return gqlClient, nil
	}

	f.NewGitHubRepo = func() (*github.Repo, error) {
		repo := &github.Repo{
			Owner: config.FromEnv("GITHUB", "REPO_OWNER"),
			Name:  config.FromEnv("GITHUB", "REPO_NAME"),
		}
		return repo, nil
	}

	return f
}
