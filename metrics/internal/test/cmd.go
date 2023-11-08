package test

import (
	"bytes"
	"context"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/spf13/cobra"
)

// ExecuteCmd uses the passed in function to create a command and execute it
func ExecuteCmd(newCmd func(*factory.Factory) *cobra.Command, args []string) (string, error) {
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
		return &export.WriterExporter{W: buf, Encoder: encoder}, nil
	}

	// Overwrite NewGitHubGraphQLClient to return canned responses (fixtures)
	// rather than querying the live GitHub GraphQL API.
	factory.NewGitHubGraphQLClient = func() (github.GraphQLClient, error) {
		return &FakeGraphQLClient{}, nil
	}

	cmd := newCmd(factory)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.ExecuteContext(ctx)

	return buf.String(), err
}
