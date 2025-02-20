package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

type CompareQuery struct {
	Repository struct {
		Name  string
		Owner struct {
			Login string
		}
		Ref struct {
			Compare struct {
				Commits struct {
					Nodes []Commit
				} `graphql:"commits(first: $first)"`
			} `graphql:"compare(headRef: $headRef)"`
		} `graphql:"ref(qualifiedName: $baseRef)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func QueryCompareRefs(gqlClient GraphQLClient, repo *Repo, base string, head string, limit int) ([]Commit, error) {
	queryVariables := map[string]interface{}{
		"owner":   githubv4.String(repo.Owner),
		"name":    githubv4.String(repo.Name),
		"baseRef": githubv4.String(base),
		"headRef": githubv4.String(head),
		"first":   githubv4.Int(limit),
	}

	var commits []Commit

	var query CompareQuery

	err := gqlClient.Query(context.Background(), &query, queryVariables)
	if err != nil {
		return nil, err
	}

	commits = append(commits, query.Repository.Ref.Compare.Commits.Nodes...)

	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits between refs")
	}

	return commits, nil
}
