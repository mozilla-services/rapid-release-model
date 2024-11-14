package io

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/mozilla-services/rapid-release-model/ciplatforms/internal/github"
)

var repoPattern = regexp.MustCompile(`^(?P<owner>[a-zA-Z0-9][a-zA-Z0-9._-]*)/(?P<name>[a-zA-Z0-9._-]+)$`)

type ServicesReader interface {
	ReadServices(filename string) ([]github.Service, map[string]*github.Repository, error)
}

type CSVServicesReader struct{}

// ReadServices loads service information from the given CSV file.
func (c CSVServicesReader) ReadServices(inputFile string) ([]github.Service, map[string]*github.Repository, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file at %s: %w", inputFile, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding CSV file at %s: %w", inputFile, err)
	}

	distinctRepos := make(map[string]*github.Repository)
	var services []github.Service

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 2 {
			return nil, nil, fmt.Errorf("invalid CSV file format: expected at least 2 columsn")
		}

		service, githubRepo := record[0], record[1]

		match := repoPattern.FindStringSubmatch(githubRepo)
		if match == nil {
			return nil, nil, fmt.Errorf("invalid GitHub repository format for %s", githubRepo)
		}

		owner := match[repoPattern.SubexpIndex("owner")]
		name := match[repoPattern.SubexpIndex("name")]

		key := fmt.Sprintf("%s/%s", owner, name)
		repo, exists := distinctRepos[key]
		if !exists {
			repo = &github.Repository{Owner: owner, Name: name}
			distinctRepos[key] = repo
		}

		services = append(services, github.Service{Name: service, Repository: repo})
	}

	log.Printf("[INFO] Read %d services (linked to %d distinct repos) from %s", len(services), len(distinctRepos), inputFile)

	return services, distinctRepos, nil
}
