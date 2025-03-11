package rest

import (
	"context"
	"log/slog"
	"net/http"

	ghrest "github.com/google/go-github/v68/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
)

// Client defines the capabilities supported by this package.
type Client interface {
	CompareCommits(ctx context.Context, owner, repo string, base, head string, opts *ghrest.ListOptions) (*ghrest.CommitsComparison, *ghrest.Response, error)
}

// GitHubRESTClient implements the Client interface and forwards calls to the
// underlying go-github client while converting model types.
type GitHubRESTClient struct {
	client *ghrest.Client
}

func NewGitHubRESTClient(httpClient *http.Client) *GitHubRESTClient {
	return &GitHubRESTClient{
		client: ghrest.NewClient(httpClient),
	}
}

// Compile-time interface assertions ensure that GitHubRESTClient implements the
// Client interface. If it fails to satisfy the interface, the compiler will
// produce an error.
var (
	_ Client = (*GitHubRESTClient)(nil)
)

// API provides access to GitHub's REST API, implementing various capabilities
// such as comparing commits.
type API struct {
	client Client
	logger *slog.Logger
}

// NewGitHubRESTAPI creates a new REST API.
func NewGitHubRESTAPI(client Client, logger *slog.Logger) *API {
	if logger == nil {
		logger = slog.Default()
	}
	return &API{client: client, logger: logger}
}

// Compile-time interface assertions ensure that API implements the required
// service interfaces. If API fails to satisfy any of these interfaces, the
// compiler will produce an error. This approach enforces interface compliance
// without requiring runtime checks.
var (
	_ github.CommitsComparisonService = (*API)(nil)
)
