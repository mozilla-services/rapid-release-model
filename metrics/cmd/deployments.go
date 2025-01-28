package cmd

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type DeploymentsOptions struct {
	Limit        int
	Environments *[]string
}

func newDeploymentsCmd(f *factory.Factory) *cobra.Command {
	opts := new(DeploymentsOptions)

	cmd := &cobra.Command{
		Use:   "deployments",
		Short: "Retrieve data about GitHub Deployments",
		Long:  "Retrieve data about GitHub Deployments",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.Limit < 1 {
				return fmt.Errorf("Limit cannot be smaller than 1.")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployments(cmd.Root().Context(), f, opts)
		},
	}
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 10, "limit for how many Deployments to fetch")

	opts.Environments = cmd.Flags().StringArray("env", nil, "multiple use for Deployment environments")

	return cmd
}

func runDeployments(ctx context.Context, f *factory.Factory, opts *DeploymentsOptions) error {
	repo, err := f.NewGitHubRepo()
	if err != nil {
		return err
	}

	gqlClient, err := f.NewGitHubGraphQLClient()
	if err != nil {
		return err
	}

	deployments, err := github.QueryDeployments(gqlClient, repo, opts.Limit, opts.Environments)
	if err != nil {
		return err
	}

	exporter, err := f.NewExporter()
	if err != nil {
		return err
	}

	return exporter.Export(deployments)
}
