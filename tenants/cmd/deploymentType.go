package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mozilla-services/rapid-release-model/tenants/internal/io"
	"github.com/spf13/cobra"
)

// deploymentTypeCmd represents the deploymentType command
var deploymentTypeCmd = &cobra.Command{
	Use:   "deploymentType",
	Short: "Parse tenant files directory for deployment type",
	Long: `Parse a local tenant files directory for deployment type
Example:
  ./tenants deploymentType -d ../../global-platform-admin/tenants`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAnalysis(cmd)
	},
}

func init() {
	rootCmd.AddCommand(deploymentTypeCmd)

	deploymentTypeCmd.Flags().StringP("directory", "d", "../tenants", "Path to the tenants directory containing YAML files")
	deploymentTypeCmd.Flags().StringP("output", "o", "deployment_type.csv", "Output file")
	deploymentTypeCmd.Flags().StringP("format", "f", "csv", "Output format: csv or json")
}

func runAnalysis(cmd *cobra.Command) error {
	directory, err := cmd.Flags().GetString("directory")
	if err != nil {
		return fmt.Errorf("failed to get directory flag: %v", err)
	}

	outputFile, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("failed to get output flag: %v", err)
	}

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return fmt.Errorf("failed to get format flag: %v", err)
	}

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
			fmt.Printf("Error reading tenant file %s: %v\n", fileName, err)
			continue
		}
		deploymentType, migrationStatus := io.ParseDeploymentTypeAndMigration(string(content))
		results = append(results, map[string]string{
			"Tenant Name":      fileName,
			"Deployment Type":  deploymentType,
			"Migration Status": migrationStatus,
		})

		if deploymentType == "argocd" {
			argocdTenants++
		}
	}

	percentageArgocd := float64(argocdTenants) / float64(totalTenants) * 100
	fmt.Printf("Total tenants: %d\n", totalTenants)
	fmt.Printf("Tenants on argocd: %d\n", argocdTenants)
	fmt.Printf("Percentage of tenants on argocd: %.2f%%\n", percentageArgocd)

	if format == "csv" {
		return io.WriteCSV(outputFile, results)
	} else if format == "json" {
		return io.WriteJSON(outputFile, results)
	} else {
		return fmt.Errorf("invalid format: %s, must be 'csv' or 'json'", format)
	}
}
