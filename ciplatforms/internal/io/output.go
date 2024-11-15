package io

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mozilla-services/rapid-release-model/ciplatforms/internal/github"
)

type ResultWriter interface {
	WriteResults(filename string, services []github.Service) error
}

type JSONResultWriter struct{}

// WriteResults saves the results to a JSON file.
func (j JSONResultWriter) WriteResults(filename string, services []github.Service) error {
	data, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

type CSVResultWriter struct{}

// WriteResults saves the results to a CSV file.
func (c CSVResultWriter) WriteResults(filename string, services []github.Service) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filename, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{"service", "repo", "circleci", "github_actions", "taskcluster", "accessible"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header to CSV file: %w", err)
	}

	// Write each service's data as a CSV row
	for _, service := range services {
		row := []string{
			service.Name,
			fmt.Sprintf("%s/%s", service.Repository.Owner, service.Repository.Name),
			fmt.Sprintf("%t", service.Repository.CircleCI),
			fmt.Sprintf("%t", service.Repository.GitHubActions),
			fmt.Sprintf("%t", service.Repository.Taskcluster),
			fmt.Sprintf("%t", service.Repository.Accessible),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing row to CSV file: %w", err)
		}
	}

	return nil
}
