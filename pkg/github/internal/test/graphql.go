package test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/mozilla-services/rapid-release-model/pkg/github/graphql"
	"github.com/shurcooL/githubv4"
)

// GraphQLQueryKeyExtra holds extra request query variables.
type GraphQLQueryKeyExtra struct {
	Environments string
}

// GraphQLQueryKey uniquely identifies a GraphQL query request.
type GraphQLQueryKey struct {
	QueryType string
	RepoOwner string
	RepoName  string
	EndCursor string
	Extra     GraphQLQueryKeyExtra
}

// GraphQLResponse represents a stored response for a GraphQL query,
// containing the raw JSON content and any associated error.
type GraphQLResponse struct {
	Content string
	Err     error
}

// FakeGraphQLClient is a mock implementation of the GitHub GraphQL API client.
type FakeGraphQLClient struct {
	mu        sync.Mutex
	responses map[GraphQLQueryKey]*GraphQLResponse
}

// NewFakeGraphQLClient initializes an in-memory fake GraphQL client.
func NewFakeGraphQLClient() *FakeGraphQLClient {
	return &FakeGraphQLClient{
		responses: make(map[GraphQLQueryKey]*GraphQLResponse),
	}
}

// RegisterResponse registers a GraphQL response for a query type.
func (c *FakeGraphQLClient) RegisterResponse(queryKey GraphQLQueryKey, response *GraphQLResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.responses[queryKey] = response
}

// Query returns the registered response or an error.
func (c *FakeGraphQLClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var repoOwner string
	if ownerVar, exists := variables["owner"]; exists {
		if ownerStr, ok := ownerVar.(githubv4.String); ok {
			repoOwner = string(ownerStr)
		}
	}

	var repoName string
	if nameVar, exists := variables["name"]; exists {
		if nameStr, ok := nameVar.(githubv4.String); ok {
			repoName = string(nameStr)
		}
	}

	var endCursor string
	if cursorVar, exists := variables["endCursor"]; exists {
		if cursorStr, ok := cursorVar.(githubv4.String); ok {
			endCursor = string(cursorStr)
		}
	}

	var envs string
	if envsVar, exists := variables["environments"]; exists {
		if envsSlice, ok := envsVar.([]githubv4.String); ok {
			var result []string
			for _, name := range envsSlice {
				result = append(result, string(name))
			}
			envs = strings.Join(result, ",")
		}
	}

	queryKey := GraphQLQueryKey{
		QueryType: fmt.Sprintf("%T", q),
		RepoOwner: repoOwner,
		RepoName:  repoName,
		EndCursor: endCursor,
		Extra:     GraphQLQueryKeyExtra{Environments: envs},
	}

	// Return registered response
	if response, exists := c.responses[queryKey]; exists {
		if response.Err != nil {
			return response.Err
		}

		if err := json.Unmarshal([]byte(response.Content), q); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %v", err)
		}
		return nil
	}

	return fmt.Errorf("no response registered for query: %v", queryKey)
}

var (
	_ graphql.Client = (*FakeGraphQLClient)(nil)
)
