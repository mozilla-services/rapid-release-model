package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mozilla-services/rapid-release-model/tenants/internal/io"
	"github.com/spf13/cobra"
)

type DeploymentTypeOptions struct {
	TenantDirectory string
	OutputFile      string
	Format          string
}

func newDeploymentTypeCmd() *cobra.Command {
	opts := new(DeploymentTypeOptions)

	cmd := &cobra.Command{
		Use:   "deploymentType",
		Short: "Parse tenant files directory for deployment type",
		Long: `Parse a local tenant files directory for deployment type
	Example:
	  ./tenants deploymentType -d ../../global-platform-admin/tenants`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Bind the directory flag
			directory, err := cmd.Flags().GetString("directory")
			if err != nil {
				return fmt.Errorf("failed to get directory flag: %w", err)
			}
			opts.TenantDirectory = directory

			// Bind the output file flag
			outputFile, err := cmd.Flags().GetString("output")
			if err != nil {
				return fmt.Errorf("failed to get output flag: %w", err)
			}
			opts.OutputFile = outputFile

			// Bind the format flag
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return fmt.Errorf("failed to get format flag: %w", err)
			}
			opts.Format = format

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalysis(opts)
		},
	}

	cmd.Flags().StringP("output", "o", "deployment_type.csv", "Output file")
	cmd.Flags().StringP("format", "f", "csv", "Output format: csv or json")

	return cmd
}

func runAnalysis(opts *DeploymentTypeOptions) error {
	directory := opts.TenantDirectory
	outputFile := opts.OutputFile
	format := opts.Format

	files, err := io.GetYamlFiles(directory)
	if err != nil {
		return fmt.Errorf("error fetching YAML files: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no YAML files found in the specified directory")
	}

	results := []map[string]string{}
	totalTenants := len(files)
	argocdTenants := 0

	for _, file := range files {
		fileName := strings.TrimSuffix(filepath.Base(file), ".yaml")
		fmt.Printf("Checking tenant: %s\n", fileName)
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Skipping tenant file %s due to an error: %v\n", fileName, err)
			continue
		}
		deploymentType, migrationStatus := io.ParseDeploymentTypeAndMigration(string(content))
		results = append(results, map[string]string{
			"TenantName":      fileName,
			"DeploymentType":  deploymentType,
			"MigrationStatus": migrationStatus,
		})

		if deploymentType == "argocd" {
			argocdTenants++
		}
	}

	percentageArgocd := float64(argocdTenants) / float64(totalTenants) * 100
	fmt.Printf("Total tenants: %d\n", totalTenants)
	fmt.Printf("Tenants on Argo CD: %d\n", argocdTenants)
	fmt.Printf("Percentage of tenants on argocd: %.2f%%\n", percentageArgocd)

	if format == "csv" {
		return io.WriteCSV(outputFile, results)
	} else if format == "json" {
		return io.WriteJSON(outputFile, results)
	} else {
		return fmt.Errorf("invalid format: %s, must be 'csv' or 'json'", format)
	}
}
