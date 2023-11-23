package test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/shurcooL/githubv4"
)

type FakeGraphQLClient struct {
	repo          *github.Repo
	wantVariables map[string]interface{}
}

func (c *FakeGraphQLClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	// Verify that the query is performed for the specified GitHub repo
	if reqOwner := string(variables["owner"].(githubv4.String)); !cmp.Equal(reqOwner, c.repo.Owner) {
		return fmt.Errorf("Repo.Owner in query variables (%v) does not match app config (%v)", reqOwner, c.repo.Owner)
	}
	if reqName := string(variables["name"].(githubv4.String)); !cmp.Equal(reqName, c.repo.Name) {
		return fmt.Errorf("Repo.Name in query variables (%v) does not match app config (%v)", reqName, c.repo.Name)
	}

	for key, want := range c.wantVariables {
		if got := variables[key]; !cmp.Equal(got, want) {
			return fmt.Errorf("unexpected value for GraphQL query variable %v\n%v", key, cmp.Diff(got, want))
		}
	}

	var key string

	// Update this for other GraphQL queries under test.
	switch v := q.(type) {
	case *github.PullRequestsQuery:
		key = "prs"
	case *github.ReleasesQuery:
		key = "releases"
	case *github.DeploymentsQuery:
		key = "deployments"
	default:
		return fmt.Errorf("unsupported query: %+v", v)
	}

	// Default filename for fixtures. We don't know the endCursor before we
	// perform the request. As a result, the filename for the first page does
	// not contain the endCursor suffix. Subsequent requests have the suffix.
	filename := "query.json"

	// The endCursor variable is set, which means we're serving the next page.
	// The `after` GraphQL parameter is set to the value of `endCursor` of the
	// previous request. Add the suffix to the fixture filenamme.
	if c := variables["endCursor"]; c != (*githubv4.String)(nil) {
		filename = fmt.Sprintf("query_%s.json", c)
	}

	jsonData, err := LoadFixture(key, filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, q)
}
