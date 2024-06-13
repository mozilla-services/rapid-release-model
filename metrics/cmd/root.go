package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mozilla-services/rapid-release-model/metrics/cmd/grafana"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/spf13/cobra"
)

type MetricsOptions struct {
	Export struct {
		Encoding string
		Filename string
	}
}

// newRootCmd creates a new base command for the metrics CLI app
func newRootCmd(f *factory.Factory) *cobra.Command {
	opts := new(MetricsOptions)

	rootCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Retrieve data for measuring software delivery performance.",
		Long:  "Retrieve data for measuring software delivery performance.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch opts.Export.Encoding {
			case "json":
				f.NewEncoder = func() (export.Encoder, error) {
					return export.NewJSONEncoder()
				}
			case "csv":
				f.NewEncoder = func() (export.Encoder, error) {
					return export.NewCSVEncoder()
				}
			case "plain":
				f.NewEncoder = func() (export.Encoder, error) {
					return export.NewPlainEncoder()
				}
			default:
				return fmt.Errorf("unsupported Export.Encoding. Please use 'json', 'csv', or 'plain'.")
			}

			if opts.Export.Filename != "" {
				f.NewExporter = func() (export.Exporter, error) {
					encoder, err := f.NewEncoder()
					if err != nil {
						return nil, err
					}
					return export.NewFileExporter(encoder, opts.Export.Filename)
				}
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&opts.Export.Encoding, "encoding", "e", "json", "export encoding")
	rootCmd.PersistentFlags().StringVarP(&opts.Export.Filename, "filename", "f", "", "export to file")

	rootCmd.AddCommand(newGitHubCmd(f))
	rootCmd.AddCommand(grafana.NewGrafanaCmd(f))

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
