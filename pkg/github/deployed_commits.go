package github

import (
	"context"
	"fmt"
	"log/slog"
)

// Options for the DeploymentsService
type DeploymentOpts struct {
	Env         string
	Sha         string
	SearchLimit int
}

// Options for the CommitsComparisonService
type CommitsOpts struct {
	Limit int
}

// Options for the QueryDeployedCommits function
type DeployedCommitsOptions struct {
	Deployment *DeploymentOpts
	Commits    *CommitsOpts
}

// QueryDeployedCommits retrieves a deployment and finds the commits deployed
// between it and the previous deployment in the same environment. It uses
// DeploymentService and CommitsComparisonService to fetch and compare commits.
func QueryDeployedCommits(
	ctx context.Context,
	repo *Repo,
	d DeploymentService, c CommitsComparisonService,
	logger *slog.Logger,
	opts *DeployedCommitsOptions,
) (*DeploymentWithCommits, error) {
	logger.Debug(
		"github.QueryDeployedCommits: querying deployed commits",
		slog.String("repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Name)),
		slog.String("github.DeploymentService", fmt.Sprintf("%T", d)),
		slog.String("github.CommitsComparisonService", fmt.Sprintf("%T", c)),
		slog.Group("deployment",
			slog.String("env", opts.Deployment.Env),
			slog.String("sha", opts.Deployment.Sha),
			slog.Int("limit", opts.Deployment.SearchLimit),
		),
		slog.Group("commits",
			slog.Int("limit", opts.Commits.Limit),
		),
	)

	deployment, prev, err := d.QueryDeployment(ctx, repo, opts.Deployment.Env, opts.Deployment.Sha, opts.Deployment.SearchLimit)
	if err != nil {
		return nil, fmt.Errorf("error querying deployments: %w", err)
	}

	logger.Debug(
		"github.QueryDeployedCommits: found deployment",
		slog.String("repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Name)),
		slog.Group("deployment",
			slog.String("ref", deployment.Ref),
			slog.String("commit.SHA", deployment.Commit.SHA),
		),
		slog.Group("previous",
			slog.String("ref", prev.Ref),
			slog.String("commit.SHA", prev.Commit.SHA),
		),
	)

	base := prev.Commit.SHA
	head := deployment.Commit.SHA

	comparison, err := c.CompareCommits(ctx, repo, base, head, opts.Commits.Limit)
	if err != nil {
		return nil, fmt.Errorf("error comparing commits for %s..%s: %w", base, head, err)
	}

	deploymentWithCommits := &DeploymentWithCommits{
		Deployment:      deployment,
		DeployedCommits: comparison.Commits,
	}

	logger.Debug(
		"github.QueryDeployedCommits: found commits between deployments",
		slog.String("repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Name)),
		slog.Int("count", len(comparison.Commits)),
		slog.Group("head",
			slog.String("ref", deployment.Ref),
			slog.String("commit.SHA", deployment.Commit.SHA),
		),
		slog.Group("base",
			slog.String("ref", prev.Ref),
			slog.String("commit.SHA", prev.Commit.SHA),
		),
	)

	return deploymentWithCommits, nil
}

// Options for the DeploymentsService
type DeploymentsOpts struct {
	Envs  *[]string
	Limit int
}

// Options for the QueryDeploymentsWithCommits function
type DeploymentWithCommitsOptions struct {
	Deployments *DeploymentsOpts
	Commits     *CommitsOpts
}

// QueryDeploymentsWithCommits fetches deployments across environments and
// determines the commits deployed between each deployment and its previous one.
// It uses DeploymentsService to fetch deployments and commit ranges, and
// CommitsComparisonService to fetch commits for the identified ranges.
func QueryDeploymentsWithCommits(
	ctx context.Context,
	repo *Repo,
	d DeploymentsService, c CommitsComparisonService,
	logger *slog.Logger,
	opts *DeploymentWithCommitsOptions,
) (map[string][]*DeploymentWithCommits, error) {
	logger.Debug(
		"github.QueryDeploymentsWithCommits: querying deployments with commits",
		slog.String("repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Name)),
		slog.String("github.DeploymentService", fmt.Sprintf("%T", d)),
		slog.String("github.CommitsComparisonService", fmt.Sprintf("%T", c)),
		slog.Group("deployments",
			slog.Any("envs", opts.Deployments.Envs),
			slog.Int("limit", opts.Deployments.Limit),
		),
		slog.Group("commits",
			slog.Int("limit", opts.Commits.Limit),
		),
	)

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

	var envCounts []any
	for env, deploys := range deploysWithCommitsByEnv {
		envCounts = append(envCounts, slog.Int(env, len(deploys)))
	}

	logger.Debug(
		"github.QueryDeploymentsWithCommits: found deployments",
		slog.String("repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Name)),
		slog.Group("deployments", envCounts...),
	)

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

			logger.Debug(
				"github.QueryDeploymentsWithCommits: found commits for deployment",
				slog.String("repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Name)),
				slog.Int("count", len(comparison.Commits)),
				slog.Group("head",
					slog.String("commit.SHA", head),
				),
				slog.Group("base",
					slog.String("commit.SHA", base),
				),
			)
			envDeployments[i].DeployedCommits = comparison.Commits
		}

		lastDeployment := envDeployments[len(envDeployments)-1]
		lastDeployment.DeployedCommits = []*Commit{lastDeployment.Commit}
	}

	return deploysWithCommitsByEnv, nil
}
