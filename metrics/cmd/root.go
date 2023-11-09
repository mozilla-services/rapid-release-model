package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/spf13/cobra"
)

// newRootCmd creates a new base command for the metrics CLI app
func newRootCmd(f *factory.Factory) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Retrieve software delivery performance metrics",
		Long:  "Retrieve software delivery performance metrics",
	}
	rootCmd.AddCommand(newGitHubCmd(f))
	return rootCmd
}

// Execute the CLI application and write errors to os.Stderr
func Execute() {
	ctx := context.Background()
	factory := factory.NewFactory(ctx)
	rootCmd := newRootCmd(factory)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// New in cobra v1.8.0. See https://github.com/spf13/cobra/pull/2044
	// Run all PersistentPreRunE hooks, so we don't have to repeat factory
	// configuration or CLI flags parsing in sub commands.
	cobra.EnableTraverseRunHooks = true
}
