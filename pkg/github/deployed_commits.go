package github

import (
	"context"
	"fmt"
)

// DeploymentWithCommits represents a deployment along with its associated deployed commits.
type DeploymentWithCommits struct {
	*Deployment
	DeployedCommits []*Commit
}

// Options for the DeploymentsService
type DeploymentsOpts struct {
	Envs  *[]string
	Limit int
}

// Options for the CommitsComparisonService
type CommitsOpts struct {
	Limit int
}

type DeployedCommitsOptions struct {
	Deployments DeploymentsOpts
	Commits     CommitsOpts
}

func QueryDeploymentsWithCommits(
	ctx context.Context,
	repo *Repo,
	opts *DeployedCommitsOptions,
	d DeploymentsService, c CommitsComparisonService,
) (map[string][]*DeploymentWithCommits, error) {

	deployments, err := d.QueryDeployments(ctx, repo, opts.Deployments.Envs, opts.Deployments.Limit)
	if err != nil {
		return nil, fmt.Errorf("error querying deployments: %w", err)
	}

	deploysWithCommitsByEnv := make(map[string][]*DeploymentWithCommits)

	for i, deployment := range deployments {
		env := deployment.LatestEnvironment
		deploysWithCommitsByEnv[env] = append(
			deploysWithCommitsByEnv[env],
			&DeploymentWithCommits{Deployment: &deployments[i]},
		)
	}

	for _, envDeployments := range deploysWithCommitsByEnv {
		for i := 0; i < len(envDeployments)-1; i++ {
			// current deployment
			head := envDeployments[i].Commit.SHA

			// previous deployment
			base := envDeployments[i+1].Commit.SHA

			comparison, err := c.CompareCommits(ctx, repo, base, head, opts.Commits.Limit)
			if err != nil {
				return nil, fmt.Errorf("error comparing commits for commits %s..%s: %w", base, head, err)
			}
			envDeployments[i].DeployedCommits = comparison.Commits
		}

		lastDeployment := envDeployments[len(envDeployments)-1]
		lastDeployment.DeployedCommits = []*Commit{lastDeployment.Commit}
	}

	return deploysWithCommitsByEnv, nil
}
