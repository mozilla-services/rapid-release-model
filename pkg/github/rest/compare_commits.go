package rest

import (
	"context"
	"fmt"

	ghrest "github.com/google/go-github/v68/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
)

// CompareCommits forwards the call to the go-github client.
func (c *GitHubRESTClient) CompareCommits(ctx context.Context, owner, repo, base, head string) (*ghrest.CommitsComparison, *ghrest.Response, error) {
	// https://docs.github.com/en/rest/commits/commits?apiVersion=2022-11-28#compare-two-commits
	// This endpoint is equivalent to running the git log BASE..HEAD command, but
	// it returns commits in a different order. The git log BASE..HEAD command
	// returns commits in reverse chronological order, whereas the API returns
	// commits in chronological order.

	// When calling this endpoint without any paging parameter (per_page or
	// page), the returned list is limited to 250 commits, and the last commit
	// in the list is the most recent of the entire comparison.
	return c.client.Repositories.CompareCommits(ctx, owner, repo, base, head, nil)
}

func (a *API) CompareCommits(ctx context.Context, repo *github.Repo, base, head string, limit int) (*github.CommitsComparison, error) {
	// TODO: limit is currently not used. Need to handle pagination in the REST API call.
	restComparison, _, err := a.client.CompareCommits(ctx, repo.Owner, repo.Name, base, head)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch commit comparison: %w", err)
	}

	if *restComparison.TotalCommits == 0 {
		return nil, fmt.Errorf("no commits between commits %s..%s", base, head)
	}

	comparison := &github.CommitsComparison{TotalCommits: restComparison.TotalCommits}

	for _, commit := range restComparison.Commits {
		comparison.Commits = append(comparison.Commits, ConvertRepositoryCommit(commit))
	}

	return comparison, nil
}
