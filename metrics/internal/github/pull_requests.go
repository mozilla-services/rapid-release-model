package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

type PullRequest struct {
	ID        string
	Number    int
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  time.Time
	MergedAt  time.Time
}

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
		} `graphql:"pullRequests(states: $states, first: $limit, after: $endCursor, orderBy: $orderBy)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// QueryPullRequests fetches information about merged PRs from the GitHub GraphQL API
func QueryPullRequests(gqlClient GraphQLClient, repo *Repo, limit int) ([]PullRequest, error) {
	queryVariables := map[string]interface{}{
		"owner":     githubv4.String(repo.Owner),
		"name":      githubv4.String(repo.Name),
		"limit":     githubv4.Int(limit),
		"endCursor": (*githubv4.String)(nil), // When paginating forwards, the cursor to continue.
		"states":    []githubv4.PullRequestState{githubv4.PullRequestStateMerged},
		"orderBy":   githubv4.IssueOrder{Field: githubv4.IssueOrderFieldUpdatedAt, Direction: githubv4.OrderDirectionDesc},
	}

	var pullRequests []PullRequest

Loop:
	for {
		var query PullRequestsQuery

		err := gqlClient.Query(context.Background(), &query, queryVariables)
		if err != nil {
			return nil, err
		}

		for _, n := range query.Repository.PullRequests.Nodes {
			pullRequests = append(pullRequests, n)
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
