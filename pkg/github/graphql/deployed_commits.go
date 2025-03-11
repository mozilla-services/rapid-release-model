package graphql

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/shurcooL/githubv4"
)

type DeployedCommitsQuery struct {
	Repository struct {
		Name  string
		Owner struct {
			Login string
		}
		Deployments struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []Deployment
		} `graphql:"deployments(first: $perPage, after: $endCursor, orderBy: $orderBy, environments: $environments)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// QueryDeployment fetches Deployments from the GitHub GraphQL API and uses that
// to search for a deployment matching the given commit SHA.
func (a *API) QueryDeployment(ctx context.Context, repo *github.Repo, env string, sha string, searchLimit int) (*github.Deployment, *github.Deployment, error) {
	a.logger.Debug(
		"graphql.QueryDeployment: querying deployment",
		slog.Group("query",
			slog.String("repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Name)),
			slog.String("env", env),
			slog.String("sha", sha),
			slog.Int("searchLimit", searchLimit),
		),
	)

	perPage := 10
	if searchLimit < perPage {
		perPage = searchLimit
	}

	queryVariables := map[string]interface{}{
		"owner":        githubv4.String(repo.Owner),
		"name":         githubv4.String(repo.Name),
		"perPage":      githubv4.Int(perPage),
		"endCursor":    (*githubv4.String)(nil),
		"orderBy":      githubv4.DeploymentOrder{Field: githubv4.DeploymentOrderFieldCreatedAt, Direction: githubv4.OrderDirectionDesc},
		"environments": []githubv4.String{githubv4.String(env)},
	}

	var deployment *github.Deployment
	var prev *github.Deployment
	var counter int

	for {
		var query DeployedCommitsQuery
		err := a.client.Query(ctx, &query, queryVariables)
		if err != nil {
			return nil, nil, err
		}

		for _, d := range query.Repository.Deployments.Nodes {
			a.logger.Debug(
				"graphql.QueryDeployment: processing deployment",
				slog.Group("deployment",
					slog.String("sha", string(d.Commit.Oid)),
					slog.Time("createdAt", d.CreatedAt),
					slog.String("ref", d.Ref.Name),
				),
			)
			counter++

			if deployment == nil {
				if string(d.Commit.Oid) == sha || string(d.Commit.AbbreviatedOid) == sha {
					// Check if this is the deployment we're looking for
					deployment = ConvertDeployment(&d)

					a.logger.Debug(
						"graphql.QueryDeployment: found matching deployment",
						slog.Group("deployment",
							slog.String("sha", string(d.Commit.Oid)),
						),
					)

					continue
				}
			}

			if deployment != nil {
				// The first deployment encountered after `deployment` is `prev`
				prev = ConvertDeployment(&d)

				a.logger.Debug(
					"graphql.QueryDeployment: found previous deployment",
					slog.Group("deployment",
						slog.String("sha", string(d.Commit.Oid)),
					),
				)

				return deployment, prev, nil // Early return since we found both
			}
		}

		if counter >= searchLimit {
			if deployment == nil {
				return nil, nil, fmt.Errorf("search limit %d reached, no deployment found for SHA %s in %s", searchLimit, sha, env)
			}
			return deployment, nil, fmt.Errorf("search limit %d reached, found deployment but no previous deployment for SHA %s in %s", searchLimit, sha, env)
		}

		if !query.Repository.Deployments.PageInfo.HasNextPage {
			if deployment == nil {
				return nil, nil, fmt.Errorf("no deployment found for SHA %s in %s", sha, env)
			}
			return deployment, nil, fmt.Errorf("found deployment but no previous deployment for SHA %s in %s", sha, env)
		}

		// Double perPage to improve efficiency but cap it at 100 to respect API limits.
		// Ensure perPage never exceeds remaining searchLimit to avoid unnecessary requests.
		//
		// Example 1: searchLimit = 50, initial perPage = 10
		//   Iteration 1 -> perPage = 10, counter = 10
		//   Iteration 2 -> perPage = 20, counter = 30
		//   Iteration 3 -> perPage = 40 (uncapped) -> 20 (capped at remaining 20), counter = 50 (stop)
		//
		// Example 2: searchLimit = 15, initial perPage = 10
		//   Iteration 1 -> perPage = 10, counter = 10
		//   Iteration 2 -> perPage = 20 (uncapped) -> 5 (capped at remaining 5), counter = 15 (stop)
		perPage *= 2
		if perPage > 100 {
			perPage = 100
		}

		remaining := searchLimit - counter
		if perPage > remaining {
			perPage = remaining
		}

		a.logger.Debug(
			"graphql.QueryDeployment: increasing perPage",
			slog.Int("newPerPage", perPage),
		)

		queryVariables["perPage"] = githubv4.Int(perPage)
		queryVariables["endCursor"] = githubv4.String(query.Repository.Deployments.PageInfo.EndCursor)
	}
}
