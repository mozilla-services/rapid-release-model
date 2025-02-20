package graphql

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/shurcooL/githubv4"
)

type compareQuery struct {
	Repository struct {
		Name  string
		Owner struct {
			Login string
		}
		Ref struct {
			Compare struct {
				Commits struct {
					Nodes      []Commit
					TotalCount int
					PageInfo   struct {
						HasNextPage bool
						EndCursor   string
					}
				} `graphql:"commits(first: $perPage, after: $endCursor)"`
			} `graphql:"compare(headRef: $headRef)"`
		} `graphql:"ref(qualifiedName: $baseRef)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (a *API) QueryCompareRefs(ctx context.Context, repo *github.Repo, base string, head string, limit int) (*github.CommitsComparison, error) {
	perPage := limit
	if limit > 100 {
		perPage = 100
	}
	queryVariables := map[string]interface{}{
		"owner":     githubv4.String(repo.Owner),
		"name":      githubv4.String(repo.Name),
		"baseRef":   githubv4.String(base),
		"headRef":   githubv4.String(head),
		"perPage":   githubv4.Int(perPage),
		"endCursor": (*githubv4.String)(nil), // When paginating forwards, the cursor to continue.
	}

	comparison := &github.CommitsComparison{}

Loop:
	for {

		var query compareQuery

		err := a.client.Query(ctx, &query, queryVariables)
		if err != nil {
			return nil, err
		}

		comparison.TotalCommits = &query.Repository.Ref.Compare.Commits.TotalCount

		if *comparison.TotalCommits == 0 {
			return nil, fmt.Errorf("no commits between refs %s..%s", base, head)
		}

		for _, c := range query.Repository.Ref.Compare.Commits.Nodes {
			comparison.Commits = append(comparison.Commits, ConvertCommit(&c))
			if len(comparison.Commits) == limit {
				break Loop
			}
		}

		if !query.Repository.Ref.Compare.Commits.PageInfo.HasNextPage {
			break
		}

		queryVariables["endCursor"] = githubv4.String(query.Repository.Ref.Compare.Commits.PageInfo.EndCursor)
	}

	return comparison, nil
}
