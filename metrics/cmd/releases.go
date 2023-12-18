package cmd

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/spf13/cobra"
)

type ReleasesOptions struct {
	Limit   int
	WithPRs bool
}

func newReleasesCmd(f *factory.Factory) *cobra.Command {
	opts := new(ReleasesOptions)

	cmd := &cobra.Command{
		Use:   "releases",
		Short: "Retrieve data about GitHub Releases",
		Long:  "Retrieve data about GitHub Releases",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.Limit < 1 {
				return fmt.Errorf("Limit cannot be smaller than 1.")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReleases(cmd.Root().Context(), f, opts)
		},
	}
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 10, "limit for how many Releases to fetch")
	cmd.Flags().BoolVar(&opts.WithPRs, "prs", false, "parse PR numbers from auto-generated release notes")

	return cmd
}

func runReleases(ctx context.Context, f *factory.Factory, opts *ReleasesOptions) error {
	repo, err := f.NewGitHubRepo()
	if err != nil {
		return err
	}

	gqlClient, err := f.NewGitHubGraphQLClient()
	if err != nil {
		return err
	}

	releases, err := github.QueryReleases(gqlClient, repo, opts.Limit)
	if err != nil {
		return err
	}

	exporter, err := f.NewExporter()
	if err != nil {
		return err
	}

	if opts.WithPRs {
		var releasesWithPRs []github.ReleaseWithPRs

		for _, release := range releases {
			releasesWithPRs = append(releasesWithPRs, *github.NewReleaseWithPRs(&release))
		}
		return exporter.Export(releasesWithPRs)
	}

	return exporter.Export(releases)
}
