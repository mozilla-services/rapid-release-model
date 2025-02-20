package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

type Deployment struct {
	Description         string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	OriginalEnvironment string
	LatestEnvironment   string
	Task                string
	State               githubv4.DeploymentState
	Commit              struct {
		AbbreviatedOid string
		Oid            string
	}
}

type DeploymentsQuery struct {
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

// QueryDeployments fetches information about Deployments from the GitHub GraphQL API
func QueryDeployments(gqlClient GraphQLClient, repo *Repo, limit int, envs *[]string) ([]Deployment, error) {
	// Values of `first` and `last` must be within 1-100. See `Node limit` in
	// GitHub's GraphQL API documentation.
	perPage := limit
	if limit > 100 {
		perPage = 100
	}

	environments := []githubv4.String{}

	for _, e := range *envs {
		environments = append(environments, githubv4.String(e))
	}

	queryVariables := map[string]interface{}{
		"owner":        githubv4.String(repo.Owner),
		"name":         githubv4.String(repo.Name),
		"perPage":      githubv4.Int(perPage),
		"endCursor":    (*githubv4.String)(nil), // When paginating forwards, the cursor to continue.
		"orderBy":      githubv4.DeploymentOrder{Field: githubv4.DeploymentOrderFieldCreatedAt, Direction: githubv4.OrderDirectionDesc},
		"environments": environments,
	}

	var deployments []Deployment

Loop:
	for {
		var query DeploymentsQuery

		err := gqlClient.Query(context.Background(), &query, queryVariables)
		if err != nil {
			return nil, err
		}

		for _, n := range query.Repository.Deployments.Nodes {
			deployments = append(deployments, n)
			if len(deployments) == limit {
				break Loop
			}
		}

		if !query.Repository.Deployments.PageInfo.HasNextPage {
			break
		}

		queryVariables["endCursor"] = githubv4.String(query.Repository.Deployments.PageInfo.EndCursor)
	}

	return deployments, nil
}
