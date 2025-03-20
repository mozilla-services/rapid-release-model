package graphql

import "github.com/mozilla-services/rapid-release-model/pkg/github"

// Convert a GraphQL API Commit parent to the unified Commit parent model.
func ConvertCommitParent(p *CommitParent) *github.CommitParent {
	return &github.CommitParent{
		SHA:            string(p.Oid),
		AbbreviatedSHA: string(p.AbbreviatedOid),
	}
}

// Convert a GraphQL API Commit to the unified Commit model.
func ConvertCommit(c *Commit) *github.Commit {
	var parents []*github.CommitParent
	for _, parent := range c.Parents.Nodes {
		parents = append(parents, ConvertCommitParent(parent))
	}

	return &github.Commit{
		SHA:            string(c.Oid),
		AbbreviatedSHA: string(c.AbbreviatedOid),
		Parents:        parents,
		AuthoredDate:   c.AuthoredDate,
		CommittedDate:  c.CommittedDate,
		Message:        c.Message,
	}
}

// Convert a GraphQL API Deployment to the unified Deployment model.
func ConvertDeployment(d *Deployment) *github.Deployment {
	return &github.Deployment{
		Description:         d.Description,
		CreatedAt:           d.CreatedAt,
		UpdatedAt:           d.UpdatedAt,
		OriginalEnvironment: d.OriginalEnvironment,
		LatestEnvironment:   d.LatestEnvironment,
		Task:                d.Task,
		State:               string(d.State),
		Ref:                 d.Ref.Name,
		Commit:              ConvertCommit(&d.Commit),
	}
}

// Convert a GraphQL API Pull Request to the unified Pull Request model
func ConvertPullRequest(p *PullRequest) *github.PullRequest {
	return &github.PullRequest{
		Number:    p.Number,
		Title:     p.Title,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		ClosedAt:  p.ClosedAt,
		MergedAt:  p.MergedAt,
	}
}

// Convert a GraphQL API Release to the unified Release model
func ConvertRelease(r *Release) *github.Release {
	return &github.Release{
		Name:         r.Name,
		TagName:      r.TagName,
		IsDraft:      r.IsDraft,
		IsLatest:     r.IsLatest,
		IsPrerelease: r.IsPrerelease,
		Description:  r.Description,
		CreatedAt:    r.CreatedAt,
		PublishedAt:  r.PublishedAt,
	}
}
