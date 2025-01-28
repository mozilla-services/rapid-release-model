package cmd

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type PullRequestsOptions struct {
	Limit int
}

func newPullRequestsCmd(f *factory.Factory) *cobra.Command {
	opts := new(PullRequestsOptions)

	cmd := &cobra.Command{
		Use:   "prs",
		Short: "Retrieve data about GitHub Pull Requests",
		Long:  "Retrieve data about GitHub Pull Requests",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.Limit < 1 {
				return fmt.Errorf("Limit cannot be smaller than 1.")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPullRequests(cmd.Root().Context(), f, opts)
		},
	}
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 10, "limit for how many PRs to fetch")
	return cmd
}

func runPullRequests(ctx context.Context, f *factory.Factory, opts *PullRequestsOptions) error {
	repo, err := f.NewGitHubRepo()
	if err != nil {
		return err
	}

	gqlClient, err := f.NewGitHubGraphQLClient()
	if err != nil {
		return err
	}

	pullRequests, err := github.QueryPullRequests(gqlClient, repo, opts.Limit)
	if err != nil {
		return err
	}

	exporter, err := f.NewExporter()
	if err != nil {
		return err
	}

	return exporter.Export(pullRequests)
}
