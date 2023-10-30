package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Prefix for application specific environment variables
const envPrefix = "RRM_METRICS_"

// Repo for which to retrieve metrics
type Repo struct {
	Owner string
	Name  string
}

// Options for the cmd
type Options struct {
	Out  io.Writer
	Repo *Repo
}

// newRootCmd creates a new base command for the metrics CLI app
func newRootCmd(w io.Writer) *cobra.Command {
	opts := &Options{
		Out: w,
		Repo: &Repo{
			Owner: os.Getenv(envPrefix + "REPO_OWNER"),
			Name:  os.Getenv(envPrefix + "REPO_NAME"),
		},
	}

	rootCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Retrieve software delivery performance metrics",
		Long:  "Retrieve software delivery performance metrics",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.Repo.Owner == "" {
				return fmt.Errorf("Repo.Owner is required. Set env var or pass flag.")
			}
			if opts.Repo.Name == "" {
				return fmt.Errorf("Repo.Name is required. Set env var or pass flag.")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRoot(opts)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&opts.Repo.Owner, "repo-owner", "o", opts.Repo.Owner, "owner of the GitHub repo")
	rootCmd.PersistentFlags().StringVarP(&opts.Repo.Name, "repo-name", "n", opts.Repo.Name, "name of the GitHub repo")

	return rootCmd
}

// runRoot performs the action for the metrics CLI command
func runRoot(opts *Options) error {
	if _, err := fmt.Fprintf(opts.Out, "Retrieving metrics for %s/%s\n", opts.Repo.Owner, opts.Repo.Name); err != nil {
		return err
	}
	return nil
}

// Execute the CLI application and write errors to os.Stderr
func Execute() {
	rootCmd := newRootCmd(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
