package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// newRootCmd creates a new root cobra command.
func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ciplatforms",
		Short: "CLI app for collecting CI platform information from GitHub.",
		Long:  "CLI app for collecting CI platform information from GitHub.",
	}
	rootCmd.AddCommand(newInfoCmd())
	return rootCmd
}

// Execute creates and executes the CLI root command.
func Execute() {
	ctx := context.Background()
	rootCmd := newRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
