package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type deploymentsConfig struct {
	*githubConfig
	withCommits  bool
	limit        int
	commitLimit  int
	environments *[]string
}

func newDeploymentsCmd(f Factory, c *githubConfig) *cobra.Command {
	config := &deploymentsConfig{githubConfig: c}

	cmd := &cobra.Command{
		Use:   "deployments",
		Short: "Retrieve data about GitHub Deployments",
		Long:  "Retrieve data about GitHub Deployments",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if config.limit < 1 {
				return fmt.Errorf("limit cannot be smaller than 1")
			}

			if cmd.Flags().Changed("commit-limit") && !config.withCommits {
				return fmt.Errorf("--commit-limit requires --commits")
			}

			if config.commitLimit < 1 {
				return fmt.Errorf("commit-limit cannot be smaller than 1")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if config.withCommits {
				return runDeploymentsWithCommits(ctx, config.graphqlAPI, config.restAPI, config)
			}

			return runDeployments(ctx, config.graphqlAPI, config)
		},
	}
	cmd.Flags().IntVarP(&config.limit, "limit", "l", 10, "maximum number of deployments to fetch")
	cmd.Flags().BoolVar(&config.withCommits, "commits", false, "include deployed commits for each deployment")
	cmd.Flags().IntVar(&config.commitLimit, "commit-limit", 250, "maximum number of commits to fetch per deployment")

	config.environments = cmd.Flags().StringArray("env", nil, "multiple use for deployment environments")

	return cmd
}

func runDeploymentsWithCommits(ctx context.Context, d github.DeploymentsService, c github.CommitsComparisonService, config *deploymentsConfig) error {
	config.logger.Debug("cmd.runDeploymentsWithCommits",
		slog.String("github.DeploymentsService", fmt.Sprintf("%T", d)),
		slog.String("github.CommitsComparisonService", fmt.Sprintf("%T", c)),
		slog.Group("config",
			slog.String("repo", fmt.Sprintf("%s/%s", config.repo.Owner, config.repo.Name)),
			slog.Any("envs", *config.environments),
			slog.Int("limit", config.limit),
			slog.Int("commitLimit", config.commitLimit),
		),
	)

	opts := &github.DeploymentWithCommitsOptions{
		Deployments: &github.DeploymentsOpts{
			Envs:  config.environments,
			Limit: config.limit,
		},
		Commits: &github.CommitsOpts{
			Limit: config.commitLimit,
		},
	}

	deploymentsByEnv, err := github.QueryDeploymentsWithCommits(ctx, config.repo, d, c, config.logger, opts)
	if err != nil {
		return fmt.Errorf("error querying deployments with commits: %w", err)
	}

	return config.exporter.Export(deploymentsByEnv)
}

func runDeployments(ctx context.Context, d github.DeploymentsService, config *deploymentsConfig) error {
	config.logger.Debug("cmd.runDeployments",
		"github.DeploymentsService", fmt.Sprintf("%T", d),
		slog.Group("config",
			slog.String("repo", fmt.Sprintf("%s/%s", config.repo.Owner, config.repo.Name)),
			slog.Any("envs", *config.environments),
			slog.Int("limit", config.limit),
		),
	)

	deployments, err := d.QueryDeployments(ctx, config.repo, config.environments, config.limit)
	if err != nil {
		return fmt.Errorf("error querying deployments: %w", err)
	}

	return config.exporter.Export(deployments)
}
