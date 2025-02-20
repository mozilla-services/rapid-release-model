package graphql

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/shurcooL/githubv4"
)

type commitHistoryQuery struct {
	Repository struct {
		Object struct {
			Commit struct {
				History struct {
					PageInfo struct {
						HasNextPage bool
						EndCursor   githubv4.String
					}
					Nodes []Commit
				} `graphql:"history(first: $perPage, after: $endCursor)"`
			} `graphql:"... on Commit"`
		} `graphql:"object(oid: $oid)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (a *API) QueryHistory(ctx context.Context, repo *github.Repo, head string, base string, limit int) ([]github.Commit, error) {
	perPage := limit
	if limit > 100 {
		perPage = 100
	}

	queryVariables := map[string]interface{}{
		"owner":     githubv4.String(repo.Owner),
		"name":      githubv4.String(repo.Name),
		"oid":       githubv4.GitObjectID(head),
		"perPage":   githubv4.Int(perPage),
		"endCursor": (*githubv4.String)(nil),
	}

	var commits []github.Commit

Loop:
	for {
		var query commitHistoryQuery
		err := a.client.Query(ctx, &query, queryVariables)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch commit history: %w", err)
		}

		for _, commit := range query.Repository.Object.Commit.History.Nodes {
			if string(commit.Oid) == base {
				break Loop
			}

			commits = append(commits, *ConvertCommit(&commit))

			if len(commits) == limit {
				return nil, fmt.Errorf("QueryHistory: reached limit of %d without finding %s", limit, base)
			}
		}

		if !query.Repository.Object.Commit.History.PageInfo.HasNextPage {
			break
		}

		queryVariables["endCursor"] = githubv4.String(query.Repository.Object.Commit.History.PageInfo.EndCursor)
	}

	return commits, nil
}
