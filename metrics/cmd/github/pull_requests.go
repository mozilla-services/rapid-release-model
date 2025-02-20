package github

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type prsConfig struct {
	*githubConfig
	limit int
}

func newPullRequestsCmd(f Factory, c *githubConfig) *cobra.Command {
	config := &prsConfig{githubConfig: c}

	cmd := &cobra.Command{
		Use:   "prs",
		Short: "Retrieve data about GitHub Pull Requests",
		Long:  "Retrieve data about GitHub Pull Requests",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if config.limit < 1 {
				return fmt.Errorf("limit cannot be smaller than 1")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runPullRequests(ctx, config.graphqlAPI, config)
		},
	}

	cmd.Flags().IntVarP(&config.limit, "limit", "l", 10, "limit for how many PRs to fetch")

	return cmd
}

func runPullRequests(ctx context.Context, p github.PullRequestsService, config *prsConfig) error {
	config.logger.Debug(
		"runPullRequests",
		"github.PullRequestsService", fmt.Sprintf("%T", p),
		"repo", fmt.Sprintf("%s/%s", config.repo.Owner, config.repo.Name),
	)

	pullRequests, err := p.QueryPullRequests(ctx, config.repo, config.limit)
	if err != nil {
		return fmt.Errorf("error querying deployments: %w", err)
	}

	return config.exporter.Export(pullRequests)
}
