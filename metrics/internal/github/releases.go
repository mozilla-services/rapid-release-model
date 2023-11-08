package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

type Release struct {
	Name         string
	TagName      string
	IsDraft      bool
	IsLatest     bool
	IsPrerelease bool
	Description  string
	CreatedAt    time.Time
	PublishedAt  time.Time
}

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
		} `graphql:"releases(first: $limit, after: $endCursor, orderBy: $orderBy)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// QueryReleases fetches information about merged PRs from the GitHub GraphQL API
func QueryReleases(gqlClient GraphQLClient, repo *Repo, limit int) ([]Release, error) {
	queryVariables := map[string]interface{}{
		"owner":     githubv4.String(repo.Owner),
		"name":      githubv4.String(repo.Name),
		"limit":     githubv4.Int(limit),
		"endCursor": (*githubv4.String)(nil), // When paginating forwards, the cursor to continue.
		"orderBy":   githubv4.ReleaseOrder{Field: githubv4.ReleaseOrderFieldCreatedAt, Direction: githubv4.OrderDirectionDesc},
	}

	var releases []Release

Loop:
	for {
		var query ReleasesQuery

		err := gqlClient.Query(context.Background(), &query, queryVariables)
		if err != nil {
			return nil, err
		}

		for _, n := range query.Repository.Releases.Nodes {
			releases = append(releases, n)
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
