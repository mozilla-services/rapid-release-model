package factory

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/grafana"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github/graphql"
	"github.com/mozilla-services/rapid-release-model/pkg/github/rest"
)

// GenericFactory provides methods for configuring logging, encoding, and
// exporting.
type GenericFactory interface {
	Logger() (*slog.Logger, error)
	ConfigureLogger(io.Writer, slog.Level) *slog.Logger

	Encoder() (export.Encoder, error)
	ConfigureEncoder(string) (export.Encoder, error)

	Exporter() (export.Exporter, error)
	ConfigureExporter(string, export.Encoder) (export.Exporter, error)
}

// GitHubFactory provides methods for managing GitHub repositories, HTTP
// clients, and API clients.
type GitHubFactory interface {
	GitHubRepo() (*github.Repo, error)
	DefaultGitHubRepo() *github.Repo
	ConfigureGitHubRepo(string, string) *github.Repo

	GitHubHTTPClient() (*http.Client, error)
	ConfigureGitHubHTTPClient() (*http.Client, error)

	GitHubRestAPI() (*rest.API, error)
	ConfigureGitHubRESTAPI(*http.Client, *slog.Logger) (*rest.API, error)

	GitHubGraphQLAPI() (*graphql.API, error)
	ConfigureGitHubGraphQLAPI(*http.Client, *slog.Logger) (*graphql.API, error)
}

// GrafanaFactory provides methods for configuring Grafana clients and
// annotation filters.
type GrafanaFactory interface {
	DefaultGrafanaAnnotationsFilter() *grafana.AnnotationsFilter

	GrafanaHubHTTPClient() (grafana.HTTPClient, error)
	ConfigureGrafanaHubHTTPClient() (grafana.HTTPClient, error)
}
