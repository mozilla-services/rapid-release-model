package test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/grafana"
	"github.com/spf13/cobra"
)

// ExecuteCmd uses the passed in function to create a command and execute it
func ExecuteCmd(newCmd func(*factory.Factory) *cobra.Command, args []string, wantReqParams *WantReqParams) (string, error) {
	ctx := context.Background()
	buf := new(bytes.Buffer)

	// Create CLI factory for the tests
	factory := factory.NewFactory(ctx)

	// Overwrite NewExporter, so that we export to buf
	factory.NewExporter = func() (export.Exporter, error) {
		encoder, err := factory.NewEncoder()
		if err != nil {
			return nil, err
		}
		return export.NewWriterExporter(buf, encoder)
	}

	// Overwrite NewGitHubGraphQLClient to return a fake client that returns
	// canned responses (fixtures) rather than sending queries to the live
	// GitHub GraphQL API.
	factory.NewGitHubGraphQLClient = func() (github.GraphQLClient, error) {
		// Create a new GitHub Repo object
		repo, err := factory.NewGitHubRepo()
		if err != nil {
			return nil, fmt.Errorf("error creating GitHub Repo object")
		}

		var reqParams *GitHubReqParams

		// Testcase did specifiy expected request parameters.
		if wantReqParams != nil {
			reqParams = wantReqParams.GitHub
		}

		return &FakeGitHubGraphQLClient{repo: repo, reqParams: reqParams}, nil
	}

	// Overwrite NewGrafanaHTTPClient to return a fake client that returns
	// canned responses (fixtures) rather than sending queries to the live
	// Grafana REST API.
	factory.NewGrafanaHTTPClient = func() (grafana.HTTPClient, error) {
		var reqParams *GrafanaReqParams

		// Testcase did specifiy expected request parameters.
		if wantReqParams != nil {
			reqParams = wantReqParams.Grafana
		}
		return &FakeGrafanaClient{reqParams: reqParams}, nil
	}

	cmd := newCmd(factory)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.ExecuteContext(ctx)

	return buf.String(), err
}
