package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/mozilla-services/rapid-release-model/metrics/cmd/github"
	"github.com/mozilla-services/rapid-release-model/metrics/cmd/grafana"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/spf13/cobra"
)

type metricsOptions struct {
	exporter struct {
		Encoding string
		Filename string
	}
	debug bool
}

type rootCmdConfig struct {
	logger   *slog.Logger
	exporter export.Exporter
}

// NewRootCmd creates a new base command for the metrics CLI app
func NewRootCmd(f *factory.DefaultFactory) *cobra.Command {
	opts := new(metricsOptions)
	config := new(rootCmdConfig)

	rootCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Retrieve data for measuring software delivery performance.",
		Long:  "Retrieve data for measuring software delivery performance.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logLevel := slog.LevelInfo
			if opts.debug {
				logLevel = slog.LevelDebug
			}
			config.logger = f.ConfigureLogger(os.Stderr, logLevel)

			encoder, err := f.ConfigureEncoder(opts.exporter.Encoding)
			if err != nil {
				return fmt.Errorf("error configuring encoder: %w", err)
			}

			exporter, err := f.ConfigureExporter(opts.exporter.Filename, encoder)
			if err != nil {
				return fmt.Errorf("error configuring exporter: %w", err)
			}
			config.exporter = exporter

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&opts.exporter.Encoding, "encoding", "e", "json", "export encoding")
	rootCmd.PersistentFlags().StringVarP(&opts.exporter.Filename, "filename", "f", "", "export to file")
	rootCmd.PersistentFlags().BoolVar(&opts.debug, "debug", false, "Enable debug logging")

	rootCmd.AddCommand(github.NewGitHubCmd(f))
	rootCmd.AddCommand(grafana.NewGrafanaCmd(f))

	return rootCmd
}

// Execute the CLI application and write errors to os.Stderr
func Execute() {
	ctx := context.Background()
	factory := factory.NewDefaultFactory(ctx)
	rootCmd := NewRootCmd(factory)
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
