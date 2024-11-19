package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

// Service links to a GitHub repository.
type Service struct {
	Name       string      `json:"name"`
	Repository *Repository `json:"repository"`
}

// Repository holds CI information for a GitHub repository.
type Repository struct {
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	CircleCI      bool   `json:"circle_ci"`
	GitHubActions bool   `json:"gh_actions"`
	Taskcluster   bool   `json:"taskcluster"`
	Accessible    bool   `json:"accessible"`
	Archived      bool   `json:"archived"`
}

// Define the GraphQL query template for batch queries for multiple repos.
const queryTemplate = `
query {
{{- range $i, $repo := . }}
  repo{{ $i }}: repository(owner: "{{ $repo.Owner }}", name: "{{ $repo.Name }}") {
    name
    owner { login }
    isArchived
    circleci: object(expression: "HEAD:.circleci/config.yml") {
      ... on Blob { id }
    }
    githubActions: object(expression: "HEAD:.github/workflows") {
      ... on Tree {
        entries { name }
      }
    }
    taskcluster: object(expression: "HEAD:.taskcluster.yml") {
      ... on Blob { id }
    }
  }
{{- end }}
}
`

// See https://docs.github.com/en/graphql/guides/forming-calls-with-graphql#the-graphql-endpoint
const githubGraphQLEndpoint = "https://api.github.com/graphql"

// CheckCIConfigInBatches dynamically generates the query for each batch and parses the response.
func CheckCIConfigInBatches(ctx context.Context, token string, repos map[string]*Repository, batchSize int) error {
	var repoSlice []*Repository
	for _, r := range repos {
		repoSlice = append(repoSlice, r)
	}
	log.Printf("[INFO] Checking CI Config for %d repos (batch size %d)", len(repoSlice), batchSize)

	for i := 0; i < len(repoSlice); i += batchSize {
		end := i + batchSize
		if end > len(repos) {
			end = len(repos)
		}
		batch := repoSlice[i:end]

		// Generate the query from the template
		query, err := buildQueryFromTemplate(batch)
		if err != nil {
			return fmt.Errorf("failed to build query: %w", err)
		}

		// Execute the batch query
		responseData, err := executeQuery(ctx, token, query)
		if err != nil {
			return fmt.Errorf("GitHub API query failed: %w", err)
		}

		if err := updateRepos(batch, responseData); err != nil {
			return fmt.Errorf("parsing results failed: %w", err)
		}
	}
	return nil
}

// buildQueryFromTemplate builds the GraphQL query for the batch using a template.
func buildQueryFromTemplate(batch []*Repository) (string, error) {
	tmpl, err := template.New("graphqlQuery").Parse(queryTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, batch); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// executeQuery sends an HTTP request with the generated GraphQL query to the GitHub GraphQL API.
func executeQuery(ctx context.Context, token string, query string) (map[string]interface{}, error) {
	reqBody := map[string]string{"query": query}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", githubGraphQLEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// See https://docs.github.com/en/graphql/guides/forming-calls-with-graphql#authenticating-with-graphql
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query failed with status code %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return response, nil
}

// updateRepos parses the response and maps results to the batch repositories
func updateRepos(batch []*Repository, data map[string]interface{}) error {
	// Ensure that the top-level "data" field exists
	dataField, ok := data["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse 'data' field in response")
	}

	for i, repo := range batch {
		alias := fmt.Sprintf("repo%d", i)

		// Retrieve the repository data for the alias
		repoData, exists := dataField[alias]
		if !exists || repoData == nil {
			// Log a warning if data for the alias is missing or set to nil (likely a 404 error)
			log.Printf("[WARNING] Data for repository %s/%s (alias %s) is missing or inaccessible", repo.Owner, repo.Name, alias)
			repo.Accessible = false
			continue
		}

		repo.Accessible = true

		repoDataMap, ok := repoData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid data format for alias %s (%s/%s): expected a map but got %T", alias, repo.Owner, repo.Name, repoData)
		}

		if archived, ok := repoDataMap["isArchived"].(bool); ok {
			repo.Archived = archived
		} else {
			log.Printf("[WARNING] 'isArchived' field missing or invalid for repository %s/%s", repo.Owner, repo.Name)
		}

		// Check if CircleCI configuration file is present
		repo.CircleCI = repoDataMap["circleci"] != nil

		// Check if Taskcluster configuration file is present
		repo.Taskcluster = repoDataMap["taskcluster"] != nil

		// Check if any GitHub Actions workflow configuration files are present
		if githubActionsData, ok := repoDataMap["githubActions"].(map[string]interface{}); ok {
			if entries, ok := githubActionsData["entries"].([]interface{}); ok {
				repo.GitHubActions = hasYmlFile(entries)
			}
		}
	}
	return nil
}

// hasYmlFile checks if any entry in the entries slice is a YAML file with a .yml or .yaml extension.
func hasYmlFile(entries []interface{}) bool {
	for _, entry := range entries {
		if entryMap, ok := entry.(map[string]interface{}); ok {
			if name, ok := entryMap["name"].(string); ok {
				if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
					return true
				}
			}
		}
	}
	return false
}
