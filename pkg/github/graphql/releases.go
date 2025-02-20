package graphql

import (
	"context"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/shurcooL/githubv4"
)

type ReleasesQuery struct {
	Repository struct {
		Name  string
		Owner struct {
			Login string
		}
		Releases struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []Release
		} `graphql:"releases(first: $perPage, after: $endCursor, orderBy: $orderBy)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// QueryReleases fetches information about merged PRs from the GitHub GraphQL API
func (a *API) QueryReleases(ctx context.Context, repo *github.Repo, limit int) ([]github.Release, error) {
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
		"orderBy":   githubv4.ReleaseOrder{Field: githubv4.ReleaseOrderFieldCreatedAt, Direction: githubv4.OrderDirectionDesc},
	}

	var releases []github.Release

Loop:
	for {
		var query ReleasesQuery

		err := a.client.Query(ctx, &query, queryVariables)
		if err != nil {
			return nil, err
		}

		for _, r := range query.Repository.Releases.Nodes {
			releases = append(releases, *ConvertRelease(&r))
			if len(releases) == limit {
				break Loop
			}
		}

		if !query.Repository.Releases.PageInfo.HasNextPage {
			break
		}

		queryVariables["endCursor"] = githubv4.String(query.Repository.Releases.PageInfo.EndCursor)
	}

	return releases, nil
}
