package test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/shurcooL/githubv4"
)

type FakeGraphQLClient struct{}

func (c *FakeGraphQLClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	var key string

	// Update this for other GraphQL queries under test.
	switch v := q.(type) {
	case *github.PullRequestsQuery:
		key = "prs"
	default:
		return fmt.Errorf("unsupported query: %+v", v)
	}

	// Default filename for fixtures (first page).
	filename := "query.json"

	// If endCursor is set, we need to serve the corresponding page instead.
	if c := variables["endCursor"]; c != (*githubv4.String)(nil) {
		filename = fmt.Sprintf("query_%s.json", c)
	}

	jsonData, err := LoadFixture(key, filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, q)
}
