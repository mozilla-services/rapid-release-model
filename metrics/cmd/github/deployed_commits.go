package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type deployedCommitsConfig struct {
	*githubConfig
	searchLimit int
	commitLimit int
	sha         string
	environment string
}

func newDeployedCommitsCmd(f Factory, c *githubConfig) *cobra.Command {
	config := &deployedCommitsConfig{githubConfig: c}

	cmd := &cobra.Command{
		Use:   "deployed-commits",
		Short: "Retrieve a deployment with its commits",
		Long:  "Retrieve a deployment with its commits",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if config.searchLimit < 1 {
				return fmt.Errorf("search-limit cannot be smaller than 1")
			}

			if config.commitLimit < 1 {
				return fmt.Errorf("commit-limit cannot be smaller than 1")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			return runDeployedCommits(ctx, config.graphqlAPI, config.restAPI, config)
		},
	}
	cmd.Flags().IntVar(&config.searchLimit, "search-limit", 10, "maximum number of deployments to search")
	cmd.Flags().IntVar(&config.commitLimit, "commit-limit", 250, "maximum number of commits to fetch for the deployment")
	cmd.Flags().StringVar(&config.environment, "env", "production", "deployment environment")
	cmd.Flags().StringVar(&config.sha, "sha", "", "git commit SHA of the deployment")

	cmd.MarkFlagRequired("sha")

	return cmd
}

func runDeployedCommits(ctx context.Context, d github.DeploymentService, c github.CommitsComparisonService, config *deployedCommitsConfig) error {
	config.logger.Debug("cmd.runDeployedCommits",
		"github.DeploymentService", fmt.Sprintf("%T", d),
		"github.CommitsComparisonService", fmt.Sprintf("%T", c),
		slog.Group("config",
			slog.String("repo", fmt.Sprintf("%s/%s", config.repo.Owner, config.repo.Name)),
			slog.String("env", config.environment),
			slog.String("sha", config.sha),
			slog.Int("searchLimit", config.searchLimit),
			slog.Int("commitLimit", config.commitLimit),
		),
	)

	opts := &github.DeployedCommitsOptions{
		Deployment: &github.DeploymentOpts{
			Env:         config.environment,
			Sha:         config.sha,
			SearchLimit: config.searchLimit,
		},
		Commits: &github.CommitsOpts{
			Limit: config.commitLimit,
		},
	}

	deployment, err := github.QueryDeployedCommits(ctx, config.repo, d, c, config.logger, opts)
	if err != nil {
		return fmt.Errorf("error querying deployed commits: %w", err)
	}

	return config.exporter.Export(deployment)
}
