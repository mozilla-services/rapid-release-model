package graphql

import (
	"context"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/shurcooL/githubv4"
)

// GraphQL query for GitHub Pull Requests
type PullRequestsQuery struct {
	Repository struct {
		Name  string
		Owner struct {
			Login string
		}
		PullRequests struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []PullRequest
		} `graphql:"pullRequests(states: $states, first: $perPage, after: $endCursor, orderBy: $orderBy)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// QueryPullRequests fetches information about merged PRs from the GitHub GraphQL API
func (a *API) QueryPullRequests(ctx context.Context, repo *github.Repo, limit int) ([]github.PullRequest, error) {
	// Values of `first` and `last` must be within 1-100. See `Node limit` in
	// GitHub's GraphQL API documentation.
	perPage := limit
	if limit > 100 {
		perPage = 100
	}

	queryVariables := map[string]interface{}{
		"owner":     githubv4.String(repo.Owner),
		"name":      githubv4.String(repo.Name),
		"perPage":   githubv4.Int(perPage),
		"endCursor": (*githubv4.String)(nil), // When paginating forwards, the cursor to continue.
		"states":    []githubv4.PullRequestState{githubv4.PullRequestStateMerged},
		"orderBy":   githubv4.IssueOrder{Field: githubv4.IssueOrderFieldUpdatedAt, Direction: githubv4.OrderDirectionDesc},
	}

	var pullRequests []github.PullRequest

Loop:
	for {
		var query PullRequestsQuery

		err := a.client.Query(ctx, &query, queryVariables)
		if err != nil {
			return nil, err
		}

		for _, p := range query.Repository.PullRequests.Nodes {
			pullRequests = append(pullRequests, *ConvertPullRequest(&p))
			if len(pullRequests) == limit {
				break Loop
			}
		}

		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}

		queryVariables["endCursor"] = githubv4.String(query.Repository.PullRequests.PageInfo.EndCursor)
	}

	return pullRequests, nil
}
