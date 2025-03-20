package rest

import (
	"context"
	"fmt"
	"log/slog"

	ghrest "github.com/google/go-github/v68/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
)

// CompareCommits forwards the call to the go-github client.
func (c *GitHubRESTClient) CompareCommits(ctx context.Context, owner, repo, base, head string, opts *ghrest.ListOptions) (*ghrest.CommitsComparison, *ghrest.Response, error) {
	// https://docs.github.com/en/rest/commits/commits?apiVersion=2022-11-28#compare-two-commits
	// This endpoint is equivalent to running the git log BASE..HEAD command, but
	// it returns commits in a different order. The git log BASE..HEAD command
	// returns commits in reverse chronological order, whereas the API returns
	// commits in chronological order.

	// When calling this endpoint without any paging parameter (per_page or
	// page), the returned list is limited to 250 commits, and the last commit
	// in the list is the most recent of the entire comparison.
	return c.client.Repositories.CompareCommits(ctx, owner, repo, base, head, opts)
}

func (a *API) CompareCommits(ctx context.Context, repo *github.Repo, base, head string, limit int) (*github.CommitsComparison, error) {
	a.logger.Debug(
		"rest.CompareCommits: comparing commits",
		slog.Group("query",
			slog.String("repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Name)),
			slog.String("base", base),
			slog.String("head", head),
			slog.Int("limit", limit),
		),
	)

	perPage := limit
	if limit > 100 {
		perPage = 100
	}

	opts := &ghrest.ListOptions{PerPage: perPage, Page: 1}

	comparison := &github.CommitsComparison{}

	for {
		a.logger.Debug(
			"rest.CompareCommits: requesting page",
			slog.Group("query",
				slog.Int("page", opts.Page),
				slog.Int("perPage", opts.PerPage),
			),
		)

		restComparison, resp, err := a.client.CompareCommits(ctx, repo.Owner, repo.Name, base, head, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch commit comparison: %w", err)
		}

		if *restComparison.TotalCommits == 0 {
			return nil, fmt.Errorf("no commits between commits %s..%s", base, head)
		}

		comparison.TotalCommits = *restComparison.TotalCommits

		a.logger.Debug(
			"rest.CompareCommits: found commits",
			slog.Int("count", len(restComparison.Commits)),
			slog.Int("total", *restComparison.TotalCommits),
			slog.Group("query",
				slog.Int("page", opts.Page),
				slog.Int("perPage", opts.PerPage),
			),
		)

		for _, commit := range restComparison.Commits {
			comparison.Commits = append(comparison.Commits, ConvertRepositoryCommit(commit))
		}

		if limit > 0 && len(comparison.Commits) >= limit {
			a.logger.Debug(
				"rest.CompareCommits: reached limit. truncating",
				slog.Int("limit", limit),
				slog.Group("query",
					slog.Int("page", opts.Page),
				),
			)
			comparison.Commits = comparison.Commits[:limit] // Truncate if needed
			break
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return comparison, nil
}
