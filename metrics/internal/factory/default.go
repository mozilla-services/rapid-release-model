package factory

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/grafana"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github/graphql"
	"github.com/mozilla-services/rapid-release-model/pkg/github/rest"
	"golang.org/x/oauth2"
)

type DefaultFactory struct {
	logger    *slog.Logger
	NewLogger func(io.Writer, slog.Level) *slog.Logger

	encoder    export.Encoder
	newEncoder func(string) (export.Encoder, error)

	exporter    export.Exporter
	NewExporter func(string, export.Encoder) (export.Exporter, error)

	githubRepo    *github.Repo
	newGitHubRepo func(string, string) *github.Repo

	githubHTTPClient    *http.Client
	newGitHubHTTPClient func() (*http.Client, error)

	githubRESTAPI       *rest.API
	NewGitHubRESTClient func(*http.Client) rest.Client
	newGitHubRESTAPI    func(rest.Client) (*rest.API, error)

	githubGraphQLAPI       *graphql.API
	NewGitHubGraphQLClient func(*http.Client) graphql.Client
	newGitHubGraphQLAPI    func(graphql.Client) (*graphql.API, error)

	grafanaHTTPClient    grafana.HTTPClient
	NewGrafanaHTTPClient func() (grafana.HTTPClient, error)
}

func NewDefaultFactory(ctx context.Context) *DefaultFactory {
	f := new(DefaultFactory)

	f.NewLogger = newLogger()
	f.newEncoder = newEncoder()
	f.NewExporter = newExporter()

	f.newGitHubRepo = newGitHubRepo()
	f.newGitHubHTTPClient = newGitHubHTTPClient(ctx)
	f.NewGitHubRESTClient = newGitHubRESTClient(ctx)
	f.newGitHubRESTAPI = newGitHubRESTAPI(ctx)
	f.NewGitHubGraphQLClient = newGitHubGraphQLClient(ctx)
	f.newGitHubGraphQLAPI = newGitHubGraphQLAPI(ctx)

	f.NewGrafanaHTTPClient = newGrafanaHTTPClient(ctx)

	return f
}

// ConfigureLogger sets the logger using the given writer and level.
func (f *DefaultFactory) ConfigureLogger(w io.Writer, l slog.Level) *slog.Logger {
	f.logger = f.NewLogger(w, l)
	return f.logger
}

// Logger returns the configured logger or an error if it is unset.
func (f *DefaultFactory) Logger() (*slog.Logger, error) {
	if f.logger == nil {
		return nil, fmt.Errorf("logger not configured")
	}
	return f.logger, nil
}

// ConfigureEncoder sets the logger using the given writer and level.
func (f *DefaultFactory) ConfigureEncoder(e string) (export.Encoder, error) {
	encoder, err := f.newEncoder(e)
	if err != nil {
		return nil, fmt.Errorf("error creating a new encoder: %w", err)
	}
	f.encoder = encoder

	return f.encoder, nil
}

// Encoder returns the configured encoder or an error if it is unset.
func (f *DefaultFactory) Encoder() (export.Encoder, error) {
	if f.encoder == nil {
		return nil, fmt.Errorf("encoder not configured")
	}
	return f.encoder, nil
}

func (f *DefaultFactory) ConfigureExporter(filename string, encoder export.Encoder) (export.Exporter, error) {
	exporter, err := f.NewExporter(filename, encoder)
	if err != nil {
		return nil, fmt.Errorf("error creating exporter: %w", err)
	}
	f.exporter = exporter

	return f.exporter, nil
}

// Exporter returns the configured exporter or an error if it is unset.
func (f *DefaultFactory) Exporter() (export.Exporter, error) {
	if f.exporter == nil {
		return nil, fmt.Errorf("exporter not configured")
	}
	return f.exporter, nil
}

// ConfigureGitHubRepo sets the repo using the given owner and name.
func (f *DefaultFactory) ConfigureGitHubRepo(owner, name string) *github.Repo {
	f.githubRepo = f.newGitHubRepo(owner, name)
	return f.githubRepo
}

// GitHubRepo returns the configured repo or an error if it is unset.
func (f *DefaultFactory) GitHubRepo() (*github.Repo, error) {
	if f.githubRepo == nil {
		return nil, fmt.Errorf("github repo not configured")
	}
	return f.githubRepo, nil
}

// DefaultGitHubRepo returns a GitHub repository using the default owner and name
// from the environment variables GITHUB_REPO_OWNER and GITHUB_REPO_NAME.
func (f *DefaultFactory) DefaultGitHubRepo() *github.Repo {
	repo := f.newGitHubRepo(
		config.ReadFromEnv("GITHUB", "REPO_OWNER"),
		config.ReadFromEnv("GITHUB", "REPO_NAME"),
	)
	return repo
}

// GitHubHTTPClient returns the configured GitHub HTTP client or an error if it
// has not been set.
func (f *DefaultFactory) GitHubHTTPClient() (*http.Client, error) {
	if f.githubHTTPClient == nil {
		return nil, fmt.Errorf("github HTTP client not configured")
	}
	return f.githubHTTPClient, nil
}

// ConfigureGitHubHTTPClient initializes and stores a GitHub HTTP client.
// Returns the client or an error if creation fails.
func (f *DefaultFactory) ConfigureGitHubHTTPClient() (*http.Client, error) {
	httpClient, err := f.newGitHubHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("error creating a GitHub HTTP client: %w", err)
	}
	f.githubHTTPClient = httpClient
	return f.githubHTTPClient, nil
}

// GitHubRestAPI returns the configured GitHub REST API or an error if it
// has not been set.
func (f *DefaultFactory) GitHubRestAPI() (*rest.API, error) {
	if f.githubRESTAPI == nil {
		return nil, fmt.Errorf("github REST API not configured")
	}
	return f.githubRESTAPI, nil
}

// ConfigureGitHubRESTAPI initializes and stores a GitHub REST API client.
// Uses the provided HTTP client to create and configure the REST API.
// Returns the configured API or an error if initialization fails.
func (f *DefaultFactory) ConfigureGitHubRESTAPI(httpClient *http.Client) (*rest.API, error) {
	client := f.NewGitHubRESTClient(httpClient)

	api, err := f.newGitHubRESTAPI(client)
	if err != nil {
		return nil, fmt.Errorf("error creating GitHub REST API: %w", err)
	}
	f.githubRESTAPI = api

	return f.githubRESTAPI, nil
}

// GitHubGraphQLAPI returns the configured GitHub GraphQL API or an error if it
// has not been set.
func (f *DefaultFactory) GitHubGraphQLAPI() (*graphql.API, error) {
	if f.githubGraphQLAPI == nil {
		return nil, fmt.Errorf("github GraphQL API not configured")
	}
	return f.githubGraphQLAPI, nil
}

// ConfigureGitHubGraphQLAPI initializes and stores a GitHub GraphQL API client.
// Uses the provided HTTP client to create and configure the GraphQL API.
// Returns the configured API or an error if initialization fails.
func (f *DefaultFactory) ConfigureGitHubGraphQLAPI(httpClient *http.Client) (*graphql.API, error) {
	client := f.NewGitHubGraphQLClient(httpClient)

	api, err := f.newGitHubGraphQLAPI(client)
	if err != nil {
		return nil, fmt.Errorf("error creating GitHub GraphQL API: %w", err)
	}
	f.githubGraphQLAPI = api

	return f.githubGraphQLAPI, nil
}

func (f *DefaultFactory) DefaultGrafanaAnnotationsFilter() *grafana.AnnotationsFilter {
	filter := &grafana.AnnotationsFilter{
		App:  config.ReadFromEnv("GRAFANA", "ANNOTATIONS", "APP"),
		From: config.ReadFromEnv("GRAFANA", "ANNOTATIONS", "FROM"),
		To:   config.ReadFromEnv("GRAFANA", "ANNOTATIONS", "TO"),
	}
	return filter
}

func (f *DefaultFactory) ConfigureGrafanaHubHTTPClient() (grafana.HTTPClient, error) {
	client, err := f.NewGrafanaHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("error creating Grafana HTTP Client: %w", err)
	}
	f.grafanaHTTPClient = client
	return f.grafanaHTTPClient, nil
}

func (f *DefaultFactory) GrafanaHubHTTPClient() (grafana.HTTPClient, error) {
	if f.grafanaHTTPClient == nil {
		return nil, fmt.Errorf("grafana HTTP client not configured")
	}
	return f.grafanaHTTPClient, nil
}

func newGitHubRepo() func(string, string) *github.Repo {
	return func(owner, name string) *github.Repo {
		repo := &github.Repo{
			Owner: owner,
			Name:  name,
		}
		return repo
	}
}

// create a func to return a new structured logger.
func newLogger() func(io.Writer, slog.Level) *slog.Logger {
	return func(w io.Writer, l slog.Level) *slog.Logger {
		logger := slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: l}))
		return logger
	}
}

// create a func to return a new export.Encoder.
func newEncoder() func(string) (export.Encoder, error) {
	return func(encoding string) (export.Encoder, error) {
		switch encoding {
		case "json":
			return export.NewJSONEncoder()
		case "csv":
			return export.NewCSVEncoder()
		case "plain":
			return export.NewPlainEncoder()
		default:
			return nil, fmt.Errorf("unsupported Export.Encoding. Please use 'json', 'csv', or 'plain'")
		}
	}
}

// create a func to return a new export.Exporter.
func newExporter() func(string, export.Encoder) (export.Exporter, error) {
	return func(filename string, encoder export.Encoder) (export.Exporter, error) {
		switch filename {
		case "":
			return export.NewWriterExporter(os.Stdout, encoder)
		default:
			return export.NewFileExporter(filename, encoder)
		}
	}
}

// create a func to return a new authenticated http.Client based on env vars.
func newGitHubHTTPClient(ctx context.Context) func() (*http.Client, error) {
	return func() (*http.Client, error) {
		token, err := config.ReadFromEnvE("GITHUB", "TOKEN")
		if err != nil {
			return nil, fmt.Errorf("error creating GitHub HTTP Client: %w", err)
		}
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

		return oauth2.NewClient(ctx, src), nil
	}
}

// create a func to return a new rest.API with the rest.Client.
func newGitHubRESTAPI(ctx context.Context) func(rest.Client) (*rest.API, error) {
	return func(client rest.Client) (*rest.API, error) {
		return rest.NewGitHubRESTAPI(client), nil
	}
}

// create a func to return a new graphql.API with the graphql.Client.
func newGitHubGraphQLAPI(ctx context.Context) func(graphql.Client) (*graphql.API, error) {
	return func(client graphql.Client) (*graphql.API, error) {
		return graphql.NewGitHubGraphQLAPI(client), nil
	}
}

// create a func to return a new rest.Client with the authenticated http.Client.
func newGitHubRESTClient(ctx context.Context) func(*http.Client) rest.Client {
	return func(c *http.Client) rest.Client {
		return rest.NewGitHubRESTClient(c)
	}
}

// create a func to return a new graphql.Client with the authenticated http.Client.
func newGitHubGraphQLClient(ctx context.Context) func(*http.Client) graphql.Client {
	return func(c *http.Client) graphql.Client {
		return graphql.NewGitHubGraphQLClient(c)
	}
}

func newGrafanaHTTPClient(ctx context.Context) func() (grafana.HTTPClient, error) {
	return func() (grafana.HTTPClient, error) {
		grafanaURL, err := config.ReadFromEnvE("GRAFANA", "SERVER_URL")
		if err != nil {
			return nil, fmt.Errorf("error creating Grafana HTTP Client: %w", err)
		}

		accessToken, err := config.ReadFromEnvE("GRAFANA", "TOKEN")
		if err != nil {
			return nil, fmt.Errorf("error creating Grafana HTTP Client: %w", err)
		}

		return grafana.NewClient(grafanaURL, accessToken)
	}
}
