package github

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type historyConfig struct {
	*githubConfig
	limit int
	base  string
	head  string
}

func newHistoryCmd(f Factory, c *githubConfig) *cobra.Command {
	config := &historyConfig{githubConfig: c}

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Retrieve commits in range",
		Long:  "Retrieve commits in range",
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
				return fmt.Errorf("git base and head commits are required")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runHistory(ctx, config.graphqlAPI, config)
		},
	}

	cmd.Flags().IntVarP(&config.limit, "limit", "l", 100, "limit for how many Commits to fetch")
	cmd.Flags().StringVar(&config.base, "base", "", "base Git commit")
	cmd.Flags().StringVar(&config.head, "head", "", "head Git commit")

	return cmd
}

func runHistory(ctx context.Context, h github.HistoryService, config *historyConfig) error {
	config.logger.Debug(
		"runHistory",
		"github.HistoryService", fmt.Sprintf("%T", h),
		"repo", fmt.Sprintf("%s/%s", config.repo.Owner, config.repo.Name),
	)

	commits, err := h.QueryHistory(ctx, config.repo, config.head, config.base, config.limit)
	if err != nil {
		return fmt.Errorf("error querying commit history: %w", err)
	}

	return config.exporter.Export(commits)
}
