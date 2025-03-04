package graphql

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/shurcooL/githubv4"
)

// API provides access to GitHub's GraphQL API,
// implementing various services such as deployments, pull requests, and releases.
type API struct {
	client Client
	logger *slog.Logger
}

// Compile-time interface assertions ensure that API implements the required service interfaces.
// If API fails to satisfy any of these interfaces, the compiler will produce an error.
// This approach enforces interface compliance without requiring runtime checks.
var (
	_ github.DeploymentsService   = (*API)(nil)
	_ github.PullRequestsService  = (*API)(nil)
	_ github.ReleasesService      = (*API)(nil)
	_ github.RefComparisonService = (*API)(nil)
	_ github.HistoryService       = (*API)(nil)
)

// Client is satisfied by the the githubv4.Client.
type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

func NewGitHubGraphQLAPI(client Client, logger *slog.Logger) *API {
	return &API{client: client, logger: logger}
}

func NewGitHubGraphQLClient(httpClient *http.Client) *githubv4.Client {
	return githubv4.NewClient(httpClient)
}
