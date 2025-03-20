package github

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type releasesConfig struct {
	*githubConfig
	limit   int
	withPRs bool
}

func newReleasesCmd(f Factory, c *githubConfig) *cobra.Command {
	config := &releasesConfig{githubConfig: c}

	cmd := &cobra.Command{
		Use:   "releases",
		Short: "Retrieve data about GitHub Releases",
		Long:  "Retrieve data about GitHub Releases",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if config.limit < 1 {
				return fmt.Errorf("limit cannot be smaller than 1")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runReleases(ctx, config.graphqlAPI, config)
		},
	}
	cmd.Flags().IntVarP(&config.limit, "limit", "l", 10, "limit for how many Releases to fetch")
	cmd.Flags().BoolVar(&config.withPRs, "prs", false, "parse PR numbers from auto-generated release notes")

	return cmd
}

func runReleases(ctx context.Context, r github.ReleasesService, config *releasesConfig) error {
	config.logger.Debug(
		"runReleases",
		"github.ReleasesService", fmt.Sprintf("%T", r),
		"repo", fmt.Sprintf("%s/%s", config.repo.Owner, config.repo.Name),
	)

	releases, err := r.QueryReleases(ctx, config.repo, config.limit)
	if err != nil {
		return fmt.Errorf("error querying releases: %w", err)
	}

	if config.withPRs {
		var releasesWithPRs []github.ReleaseWithPRs

		for _, release := range releases {
			r := release // Create a copy to avoid referencing the same loop variable memory.
			releasesWithPRs = append(releasesWithPRs, *github.NewReleaseWithPRs(&r))
		}

		return config.exporter.Export(releasesWithPRs)
	}

	return config.exporter.Export(releases)

}
