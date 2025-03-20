package github

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type compareConfig struct {
	*githubConfig
	limit int
	base  string
	head  string
}

func newCompareRefsCmd(f Factory, c *githubConfig) *cobra.Command {
	config := &compareConfig{githubConfig: c}

	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Retrieve commits between Git refs",
		Long:  "Retrieve commits between Git refs",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if config.limit < 1 {
				return fmt.Errorf("limit cannot be smaller than 1")
			}

			// TODO: add CSV encoding support for this command
			if encoding, err := cmd.Flags().GetString("encoding"); err != nil {
				return fmt.Errorf("failed to retrieve 'encoding' flag: %w", err)
			} else if encoding != "json" {
				return fmt.Errorf("unsupported Export.Encoding. Please use 'json'")
			}

			if config.base == "" || config.head == "" {
				return fmt.Errorf("git Ref Base and Head are required")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runCompareRefs(ctx, config.graphqlAPI, config)
		},
	}
	cmd.Flags().IntVarP(&config.limit, "limit", "l", 100, "limit for how many Commits to fetch")
	cmd.Flags().StringVar(&config.base, "base", "", "base Git ref")
	cmd.Flags().StringVar(&config.head, "head", "", "head Git ref")

	return cmd
}

func runCompareRefs(ctx context.Context, r github.RefComparisonService, config *compareConfig) error {
	config.logger.Debug(
		"runCompareRefs",
		"github.RefComparisonService", fmt.Sprintf("%T", r),
		"repo", fmt.Sprintf("%s/%s", config.repo.Owner, config.repo.Name),
	)

	comparison, err := r.QueryCompareRefs(ctx, config.repo, config.base, config.head, config.limit)
	if err != nil {
		return fmt.Errorf("error querying ref comparison: %w", err)
	}

	return config.exporter.Export(comparison)
}
