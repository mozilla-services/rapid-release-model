package cmd

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/spf13/cobra"
)

type CompareOptions struct {
	Limit int
	Base  string
	Head  string
}

func newCompareCmd(f *factory.Factory) *cobra.Command {
	opts := new(CompareOptions)

	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Retrieve commits between Git refs",
		Long:  "Retrieve commits between Git refs",
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
				return fmt.Errorf("git Ref Base and Head are required")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(cmd.Root().Context(), f, opts)
		},
	}
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 10, "limit for how many Commits to fetch")
	cmd.Flags().StringVar(&opts.Base, "base", "", "base Git ref")
	cmd.Flags().StringVar(&opts.Head, "head", "", "head Git ref")

	return cmd
}

func runCompare(ctx context.Context, f *factory.Factory, opts *CompareOptions) error {
	repo, err := f.NewGitHubRepo()
	if err != nil {
		return err
	}

	gqlClient, err := f.NewGitHubGraphQLClient()
	if err != nil {
		return err
	}

	commits, err := github.QueryCompareRefs(gqlClient, repo, opts.Base, opts.Head, opts.Limit)
	if err != nil {
		return err
	}

	exporter, err := f.NewExporter()
	if err != nil {
		return err
	}

	return exporter.Export(commits)
}
