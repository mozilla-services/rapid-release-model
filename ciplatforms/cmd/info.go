package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mozilla-services/rapid-release-model/ciplatforms/internal/github"
	"github.com/mozilla-services/rapid-release-model/ciplatforms/internal/io"
	"github.com/spf13/cobra"
)

const githubTokenEnvKey = "CIPLATFORMS_GITHUB_API_TOKEN"

// infoOptions holds options for the CLI command
type infoOptions struct {
	inputFile      string
	outputFile     string
	githubAPIToken string
	timeout        time.Duration
	batchSize      int

	// set in command PreRunE
	servicesReader io.ServicesReader
	resultWriter   io.ResultWriter
}

// newInfoCmd creates a new info CLI command
func newInfoCmd() *cobra.Command {
	opts := new(infoOptions)

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Collect CI platform information from GitHub.",
		Long:  "Collect CI platform information from GitHub.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.githubAPIToken == "" {
				val, ok := os.LookupEnv(githubTokenEnvKey)
				if !ok {
					return fmt.Errorf("GitHub API token required. Pass --gh-token or set %s", githubTokenEnvKey)
				}
				opts.githubAPIToken = val
			}

			switch ext := filepath.Ext(opts.inputFile); ext {
			case ".csv":
				opts.servicesReader = io.CSVServicesReader{}
			default:
				return fmt.Errorf("unsupported file extension: %s", ext)
			}

			switch ext := filepath.Ext(opts.outputFile); ext {
			case ".json":
				opts.resultWriter = io.JSONResultWriter{}
			case ".csv":
				opts.resultWriter = io.CSVResultWriter{}
			default:
				return fmt.Errorf("unsupported file extension: %s", ext)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(cmd.Root().Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.inputFile, "input", "i", "repos.csv", "input file")
	cmd.Flags().StringVarP(&opts.outputFile, "output", "o", "repos_ciplatforms.csv", "output file")
	cmd.Flags().StringVarP(&opts.githubAPIToken, "gh-token", "t", "", "GitHub API token")

	cmd.Flags().DurationVar(&opts.timeout, "timeout", 10*time.Second, "timeout for GitHub API requests")
	cmd.Flags().IntVar(&opts.batchSize, "batch-size", 50, "number of repositories to process in each batch")
	return cmd
}

func runInfo(ctx context.Context, opts *infoOptions) error {
	// Load services from the given input file.
	services, repos, err := opts.servicesReader.ReadServices(opts.inputFile)
	if err != nil {
		return fmt.Errorf("error loading services from file: %w", err)
	}

	// Ensure any operations that use timeoutCtx are automatically canceled
	// after 10 seconds. This includes long running HTTP requests.
	timeoutCtx, cancel := context.WithTimeout(ctx, opts.timeout)
	defer cancel()

	// Check CI Platform config files for each GitHub repository in batches.
	if err := github.CheckCIConfigInBatches(timeoutCtx, opts.githubAPIToken, repos, opts.batchSize); err != nil {
		return fmt.Errorf("error checking CI configs: %w", err)
	}

	// Write the results to the specified file.
	if err := opts.resultWriter.WriteResults(opts.outputFile, services); err != nil {
		return fmt.Errorf("failed to save results: %w", err)
	}
	log.Printf("[INFO] Results saved to %s\n", opts.outputFile)

	return nil
}
