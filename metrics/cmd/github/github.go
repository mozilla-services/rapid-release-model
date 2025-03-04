package github

import (
	"fmt"
	"log/slog"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github/graphql"
	"github.com/mozilla-services/rapid-release-model/pkg/github/rest"
	"github.com/spf13/cobra"
)

type Factory interface {
	factory.GenericFactory
	factory.GitHubFactory
}

type githubConfig struct {
	logger     *slog.Logger
	exporter   export.Exporter
	repo       *github.Repo
	graphqlAPI *graphql.API
	restAPI    *rest.API
}

func (c *githubConfig) configureAPIs(f Factory, logger *slog.Logger) error {
	var err error

	httpClient, err := f.ConfigureGitHubHTTPClient()
	if err != nil {
		return fmt.Errorf("error initializing GitHub HTTP client: %w", err)
	}

	if c.restAPI, err = f.ConfigureGitHubRESTAPI(httpClient, logger); err != nil {
		return fmt.Errorf("error initializing GitHub REST API: %w", err)
	}

	if c.graphqlAPI, err = f.ConfigureGitHubGraphQLAPI(httpClient, logger); err != nil {
		return fmt.Errorf("error initializing GitHub GraphQL API: %w", err)
	}

	return nil
}

func NewGitHubCmd(f Factory) *cobra.Command {
	config := &githubConfig{repo: f.DefaultGitHubRepo()}

	cmd := &cobra.Command{
		Use:   "github",
		Short: "Retrieve metrics from GitHub",
		Long:  "Retrieve metrics from GitHub",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger, err := f.Logger()
			if err != nil {
				return fmt.Errorf("error retrieving logger: %w", err)
			}
			config.logger = logger

			exporter, err := f.Exporter()
			if err != nil {
				return fmt.Errorf("error retrieving exporter: %w", err)
			}
			config.exporter = exporter

			if config.repo.Owner == "" || config.repo.Name == "" {
				return fmt.Errorf("repo.Owner and repo.Name are required. Set env vars or pass flags")
			}
			f.ConfigureGitHubRepo(config.repo.Owner, config.repo.Name)

			if err := config.configureAPIs(f, config.logger); err != nil {
				return fmt.Errorf("error configuring GitHub APIs: %w", err)
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&config.repo.Owner, "repo-owner", "o", config.repo.Owner, "owner of the GitHub repo")
	cmd.PersistentFlags().StringVarP(&config.repo.Name, "repo-name", "n", config.repo.Name, "name of the GitHub repo")

	cmd.AddCommand(newPullRequestsCmd(f, config))
	cmd.AddCommand(newReleasesCmd(f, config))
	cmd.AddCommand(newDeploymentsCmd(f, config))
	cmd.AddCommand(newCompareRefsCmd(f, config))
	cmd.AddCommand(newHistoryCmd(f, config))

	return cmd
}
