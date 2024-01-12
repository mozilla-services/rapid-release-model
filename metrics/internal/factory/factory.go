package factory

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/grafana"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Factory is used to create dependencies for the CLI application
type Factory struct {
	NewEncoder                  func() (export.Encoder, error)
	NewExporter                 func() (export.Exporter, error)
	NewGitHubHTTPClient         func() (*http.Client, error)
	NewGitHubGraphQLClient      func() (github.GraphQLClient, error)
	NewGitHubRepo               func() (*github.Repo, error)
	NewGrafanaHTTPClient        func() (grafana.HTTPClient, error)
	NewGrafanaAnnotationsFilter func() (*grafana.AnnotationsFilter, error)
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
		return export.NewWriterExporter(os.Stdout, encoder)
	}

	f.NewGitHubHTTPClient = func() (*http.Client, error) {
		token, err := config.ReadFromEnvE("GITHUB", "TOKEN")
		if err != nil {
			return nil, fmt.Errorf("Error creating GitHub HTTP Client: %w", err)
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
			Owner: config.ReadFromEnv("GITHUB", "REPO_OWNER"),
			Name:  config.ReadFromEnv("GITHUB", "REPO_NAME"),
		}
		return repo, nil
	}

	f.NewGrafanaHTTPClient = func() (grafana.HTTPClient, error) {
		grafanaURL, err := config.ReadFromEnvE("GRAFANA", "SERVER_URL")
		if err != nil {
			return nil, fmt.Errorf("Error creating Grafana HTTP Client: %w", err)
		}

		accessToken, err := config.ReadFromEnvE("GRAFANA", "TOKEN")
		if err != nil {
			return nil, fmt.Errorf("Error creating Grafana HTTP Client: %w", err)
		}

		return grafana.NewClient(grafanaURL, accessToken)
	}

	f.NewGrafanaAnnotationsFilter = func() (*grafana.AnnotationsFilter, error) {
		repo := &grafana.AnnotationsFilter{
			App:  config.ReadFromEnv("GRAFANA", "ANNOTATIONS", "APP"),
			From: config.ReadFromEnv("GRAFANA", "ANNOTATIONS", "FROM"),
			To:   config.ReadFromEnv("GRAFANA", "ANNOTATIONS", "TO"),
		}
		return repo, nil
	}

	return f
}
