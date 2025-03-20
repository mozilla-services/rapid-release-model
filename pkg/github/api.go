package github

import (
	"context"
)

// PullRequestsService provides access to GitHub Pull Request functionality.
type PullRequestsService interface {
	QueryPullRequests(ctx context.Context, repo *Repo, limit int) ([]PullRequest, error)
}

// DeploymentsService provides access to GitHub Deployment functionality.
type DeploymentsService interface {
	QueryDeployments(ctx context.Context, repo *Repo, envs *[]string, limit int) ([]Deployment, error)
}

// DeploymentService provides access to GitHub Deployment functionality.
type DeploymentService interface {
	QueryDeployment(ctx context.Context, repo *Repo, env string, sha string, searchLimit int) (*Deployment, *Deployment, error)
}

// ReleasesService provides access to GitHub Release functionality.
type ReleasesService interface {
	QueryReleases(ctx context.Context, repo *Repo, limit int) ([]Release, error)
}

// CommitsComparisonService provides commit comparison functionality.
type CommitsComparisonService interface {
	CompareCommits(ctx context.Context, repo *Repo, base string, head string, limit int) (*CommitsComparison, error)
}

// RefComparisonService provides Git reference comparison functionality.
type RefComparisonService interface {
	QueryCompareRefs(ctx context.Context, repo *Repo, base string, head string, limit int) (*CommitsComparison, error)
}

// HistoryService provides commit history functionality.
type HistoryService interface {
	QueryHistory(ctx context.Context, repo *Repo, head string, base string, limit int) ([]Commit, error)
}
