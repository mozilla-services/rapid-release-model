package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

type CommitHistoryQuery struct {
	Repository struct {
		Object struct {
			Commit struct {
				History struct {
					PageInfo struct {
						HasNextPage bool
						EndCursor   githubv4.String
					}
					Nodes []Commit
				} `graphql:"history(first: $first, after: $endCursor)"`
			} `graphql:"... on Commit"`
		} `graphql:"object(oid: $oid)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func QueryHistory(gqlClient GraphQLClient, repo *Repo, head string, base string, limit int) ([]Commit, error) {
	var query CommitHistoryQuery
	queryVariables := map[string]interface{}{
		"owner":     githubv4.String(repo.Owner),
		"name":      githubv4.String(repo.Name),
		"oid":       githubv4.GitObjectID(head),
		"first":     githubv4.Int(limit),
		"endCursor": (*githubv4.String)(nil),
	}

	var commits []Commit

Loop:
	for {
		err := gqlClient.Query(context.Background(), &query, queryVariables)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch commit history: %w", err)
		}

		for _, commit := range query.Repository.Object.Commit.History.Nodes {
			// Stop if the base commitSHA is reached
			if string(commit.Oid) == base {
				return commits, nil
			}

			commits = append(commits, commit)

			if len(commits) == limit {
				break Loop
			}

		}

		if !query.Repository.Object.Commit.History.PageInfo.HasNextPage {
			break
		}

		queryVariables["endCursor"] = githubv4.String(query.Repository.Object.Commit.History.PageInfo.EndCursor)
	}

	return commits, nil
}
