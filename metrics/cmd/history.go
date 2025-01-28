package cmd

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type HistoryOptions struct {
	Limit int
	Base  string
	Head  string
}

func newHistoryCmd(f *factory.Factory) *cobra.Command {
	opts := new(HistoryOptions)

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Retrieve commits in range",
		Long:  "Retrieve commits in range",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.Limit < 1 {
				return fmt.Errorf("limit cannot be smaller than 1")
			}
			// TODO: add CSV encoding support for this command
			if encoding, err := cmd.Flags().GetString("encoding"); err != nil {
				return fmt.Errorf("failed to retrieve 'encoding' flag: %w", err)
			} else if encoding != "json" {
				return fmt.Errorf("unsupported Export.Encoding. Please use 'json'")
			}
			if opts.Base == "" || opts.Head == "" {
				return fmt.Errorf("git base and head commits are required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHistory(cmd.Root().Context(), f, opts)
		},
	}
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 100, "limit for how many Commits to fetch")
	cmd.Flags().StringVar(&opts.Base, "base", "", "base Git commit")
	cmd.Flags().StringVar(&opts.Head, "head", "", "head Git commit")

	return cmd
}

func runHistory(ctx context.Context, f *factory.Factory, opts *HistoryOptions) error {
	repo, err := f.NewGitHubRepo()
	if err != nil {
		return err
	}

	gqlClient, err := f.NewGitHubGraphQLClient()
	if err != nil {
		return err
	}

	commits, err := github.QueryHistory(gqlClient, repo, opts.Head, opts.Base, opts.Limit)
	if err != nil {
		return err
	}

	exporter, err := f.NewExporter()
	if err != nil {
		return err
	}

	return exporter.Export(commits)
}
