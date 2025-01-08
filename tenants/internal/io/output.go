package io

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

func WriteCSV(outputFile string, results []map[string]string) error {
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer output.Close()

	csvWriter := csv.NewWriter(output)
	defer csvWriter.Flush()

	// Define the explicit order of fields
	headers := []string{"Tenant Name", "Deployment Type", "Migration Status"}
	csvWriter.Write(headers)

	// Write rows
	for _, result := range results {
		row := []string{
			result["Tenant Name"],
			result["Deployment Type"],
			result["Migration Status"],
		}
		csvWriter.Write(row)
	}

	return nil
}

type Result struct {
	TenantName      string `json:"Tenant Name"`
	DeploymentType  string `json:"Deployment Type"`
	MigrationStatus string `json:"Migration Status"`
}

func WriteJSON(outputFile string, results []map[string]string) error {
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer output.Close()

	// Convert map results to structured results
	structResults := []Result{}
	for _, result := range results {
		structResults = append(structResults, Result{
			TenantName:      result["Tenant Name"],
			DeploymentType:  result["Deployment Type"],
			MigrationStatus: result["Migration Status"],
		})
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(structResults); err != nil {
		return fmt.Errorf("error writing JSON output: %v", err)
	}

	return nil
}
