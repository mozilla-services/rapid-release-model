package test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/grafana"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github/graphql"
	"github.com/spf13/cobra"
)

// noopTransport is an http.RoundTripper that blocks all outgoing HTTP requests.
type noopTransport struct{}

// RoundTrip prevents any HTTP requests from being sent.
func (n *noopTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("outgoing HTTP requests are blocked in tests")
}

// create a func to return a new structured logger.
func newLogger(w io.Writer) func(io.Writer, slog.Level) *slog.Logger {
	return func(_ io.Writer, l slog.Level) *slog.Logger {
		logger := slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: l}))
		return logger
	}
}

// create a func to return a new export.Exporter.
func newExporter(w io.Writer) func(string, export.Encoder) (export.Exporter, error) {
	return func(filename string, encoder export.Encoder) (export.Exporter, error) {
		switch filename {
		case "":
			return export.NewWriterExporter(w, encoder)
		default:
			return export.NewFileExporter(filename, encoder)
		}
	}
}

// create a func to return a new graphql.Client with the authenticated http.Client.
func newGitHubGraphQLClient(wantReqParams *WantReqParams) func(*http.Client) graphql.Client {
	return func(c *http.Client) graphql.Client {
		repo := &github.Repo{Owner: "hackebrot", Name: "turtle"}

		var reqParams *GitHubReqParams

		// Testcase did specifiy expected request parameters.
		if wantReqParams != nil {
			reqParams = wantReqParams.GitHub
		}

		return &FakeGitHubGraphQLClient{repo: repo, reqParams: reqParams}
	}
}

// ExecuteCmd uses the passed in function to create a command and execute it
func ExecuteCmd(newCmd func(factory.Factory) *cobra.Command, args []string, wantReqParams *WantReqParams) (string, string, error) {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	logbuf := new(bytes.Buffer)

	// Create CLI f for the tests
	f := factory.NewDefaultFactory(ctx)

	// Overwrite NewExporter, so that we export to buf
	f.NewExporter = newExporter(buf)

	// Overwrite NewLogger, so that we log to logbuf
	f.NewLogger = newLogger(logbuf)

	// Override NewGitHubHTTPClient to return an http.Client that prevents
	// outgoing GitHub API requests while satisfying the factory. Without this,
	// the factory fails to configure the GitHub GraphQL and REST API clients due
	// to a missing GitHub token environment variable. Fake API clients in tests
	// do not use this http.Client.
	f.NewGitHubHTTPClient = func() (*http.Client, error) {
		return &http.Client{Transport: &noopTransport{}}, nil
	}

	// Overwrite NewGitHubGraphQLClient to return a fake client that returns
	// canned responses (fixtures) rather than sending queries to the live
	// GitHub GraphQL API.
	f.NewGitHubGraphQLClient = newGitHubGraphQLClient(wantReqParams)

	// Overwrite NewGrafanaHTTPClient to return a fake client that returns
	// canned responses (fixtures) rather than sending queries to the live
	// Grafana REST API.
	f.NewGrafanaHTTPClient = func() (grafana.HTTPClient, error) {
		var reqParams *GrafanaReqParams

		// Testcase did specifiy expected request parameters.
		if wantReqParams != nil {
			reqParams = wantReqParams.Grafana
		}
		return &FakeGrafanaClient{reqParams: reqParams}, nil
	}

	cmd := newCmd(f)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.ExecuteContext(ctx)

	return buf.String(), logbuf.String(), err
}
