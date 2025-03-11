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
	ConfigureLogger(io.Writer, slog.Level)

	Encoder() (export.Encoder, error)
	ConfigureEncoder(string) error

	Exporter() (export.Exporter, error)
	ConfigureExporter(string) error
}

// GitHubFactory provides methods for managing GitHub repositories, HTTP
// clients, and API clients.
type GitHubFactory interface {
	GitHubRepo() (*github.Repo, error)
	DefaultGitHubRepo() *github.Repo
	ConfigureGitHubRepo(string, string)

	GitHubHTTPClient() (*http.Client, error)
	ConfigureGitHubHTTPClient() error

	GitHubRestAPI() (*rest.API, error)
	ConfigureGitHubRESTAPI() error

	GitHubGraphQLAPI() (*graphql.API, error)
	ConfigureGitHubGraphQLAPI() error
}

// GrafanaFactory provides methods for configuring Grafana clients and
// annotation filters.
type GrafanaFactory interface {
	DefaultGrafanaAnnotationsFilter() *grafana.AnnotationsFilter

	GrafanaHubHTTPClient() (grafana.HTTPClient, error)
	ConfigureGrafanaHubHTTPClient() error
}
