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
	headers := []string{"TenantName", "DeploymentType", "MigrationStatus"}
	err = csvWriter.Write(headers)
	if err != nil {
		fmt.Printf("Error writing CSV Headers")
	}

	// Write rows
	for _, result := range results {
		row := []string{
			result["TenantName"],
			result["DeploymentType"],
			result["MigrationStatus"],
		}
		err = csvWriter.Write(row)
		if err != nil {
			fmt.Printf("Error writing CSV Rows")
		}
	}

	return nil
}

type Result struct {
	TenantName      string
	DeploymentType  string
	MigrationStatus string
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
			TenantName:      result["TenantName"],
			DeploymentType:  result["DeploymentType"],
			MigrationStatus: result["MigrationStatus"],
		})
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(structResults); err != nil {
		return fmt.Errorf("error writing JSON output: %v", err)
	}

	return nil
}
