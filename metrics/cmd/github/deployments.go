package github

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type deploymentsConfig struct {
	*githubConfig
	withCommits  bool
	limit        int
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
	cmd.Flags().IntVarP(&config.limit, "limit", "l", 10, "limit for how many Deployments to fetch")
	cmd.Flags().BoolVar(&config.withCommits, "commits", false, "get deployed commits for each deployment")

	config.environments = cmd.Flags().StringArray("env", nil, "multiple use for Deployment environments")

	return cmd
}

func runDeploymentsWithCommits(ctx context.Context, d github.DeploymentsService, c github.CommitsComparisonService, config *deploymentsConfig) error {
	config.logger.Debug(
		"runDeploymentsWithCommits",
		"github.DeploymentsService", fmt.Sprintf("%T", d),
		"github.CommitsComparisonService", fmt.Sprintf("%T", c),
		"repo", fmt.Sprintf("%s/%s", config.repo.Owner, config.repo.Name),
		"envs", *config.environments,
	)

	opts := &github.DeployedCommitsOptions{
		Deployments: github.DeploymentsOpts{
			Envs:  config.environments,
			Limit: config.limit,
		},
		Commits: github.CommitsOpts{
			Limit: 250,
		},
	}

	deploymentsByEnv, err := github.QueryDeploymentsWithCommits(ctx, config.repo, opts, d, c)
	if err != nil {
		return fmt.Errorf("error querying deployments with commits: %w", err)
	}

	return config.exporter.Export(deploymentsByEnv)
}

func runDeployments(ctx context.Context, d github.DeploymentsService, config *deploymentsConfig) error {
	config.logger.Debug(
		"runDeployments",
		"github.DeploymentsService", fmt.Sprintf("%T", d),
		"repo", fmt.Sprintf("%s/%s", config.repo.Owner, config.repo.Name),
		"envs", *config.environments,
	)

	deployments, err := d.QueryDeployments(ctx, config.repo, config.environments, config.limit)
	if err != nil {
		return fmt.Errorf("error querying deployments: %w", err)
	}

	return config.exporter.Export(deployments)
}
