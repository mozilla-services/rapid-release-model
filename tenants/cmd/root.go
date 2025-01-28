package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// newRootCmd creates a new base command for the metrics CLI app
func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tenants",
		Short: "CLI app for parsing GCPv2 tenant files for cataloguing",
		Long:  "CLI app for parsing GCPv2 tenant files for cataloguing",
	}

	rootCmd.PersistentFlags().StringP("directory", "d", "", "Path to the tenants directory containing YAML files")
	rootCmd.MarkPersistentFlagRequired("directory")

	rootCmd.AddCommand(newDeploymentTypeCmd())

	return rootCmd
}

// Execute the CLI application and write errors to os.Stderr
func Execute() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
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
