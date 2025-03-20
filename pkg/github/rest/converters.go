package rest

import (
	"time"

	ghrest "github.com/google/go-github/v68/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
)

func abbreviateSHA(s string) string {
	if len(s) > 7 {
		return s[:7]
	}
	return s
}

func ConvertCommitParent(p *ghrest.Commit) *github.CommitParent {
	sha := p.GetSHA()

	return &github.CommitParent{
		SHA:            sha,
		AbbreviatedSHA: abbreviateSHA(sha),
	}
}

// Convert REST API Commit to unified Commit
func ConvertCommit(c *ghrest.Commit) *github.Commit {
	sha := c.GetSHA()

	return &github.Commit{
		SHA:            sha,
		AbbreviatedSHA: abbreviateSHA(sha),
	}
}

// Convert REST API Commit to unified Commit
func ConvertRepositoryCommit(c *ghrest.RepositoryCommit) *github.Commit {
	var parents []*github.CommitParent
	for _, parent := range c.Parents {
		parents = append(parents, ConvertCommitParent(parent))
	}

	sha := c.GetSHA()

	return &github.Commit{
		SHA:            sha,
		AbbreviatedSHA: abbreviateSHA(sha),
		Parents:        parents,
		AuthoredDate:   safeGetCommitDate(c.Commit.Author),
		CommittedDate:  safeGetCommitDate(c.Commit.Committer),
		Message:        c.Commit.GetMessage(),
	}
}

// Helper function to safely extract commit date
func safeGetCommitDate(author *ghrest.CommitAuthor) time.Time {
	if author == nil {
		return time.Time{}
	}
	return author.GetDate().Time
}
