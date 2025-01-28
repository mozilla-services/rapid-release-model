package github

import (
	"context"
)

// Represents a GitHub repo
type Repo struct {
	Owner string
	Name  string
}

// GraphQLClient is satisfied by the the githubv4.Client.
type GraphQLClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}
