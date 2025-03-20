package graphql

import (
	"time"

	"github.com/shurcooL/githubv4"
)

// CommitParent represents a GitHub GraphQL API Commit without `parents`.
// See https://docs.github.com/en/graphql/reference/objects#commit
type CommitParent struct {
	AbbreviatedOid githubv4.GitObjectID
	Oid            githubv4.GitObjectID
}

// Commit represents a GitHub GraphQL API Commit.
// See https://docs.github.com/en/graphql/reference/objects#commit
type Commit struct {
	AbbreviatedOid githubv4.GitObjectID
	Oid            githubv4.GitObjectID
	Parents        struct {
		// Using *Commit here results in an error
		Nodes []*CommitParent
	} `graphql:"parents(first: 2)"`
	AuthoredDate  time.Time
	CommittedDate time.Time
	Message       string
}

// Deployment represents a GitHub GraphQL API Deployment.
// See https://docs.github.com/en/graphql/reference/objects#deployment
type Deployment struct {
	Description         string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	OriginalEnvironment string
	LatestEnvironment   string
	Task                string
	State               githubv4.DeploymentState
	Commit              Commit
	Ref                 struct {
		Name string
	}
}

// PullRequest represents a GitHub GraphQL API Pull Request.
// See https://docs.github.com/en/graphql/reference/objects#pullrequest
type PullRequest struct {
	Number    int
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  time.Time
	MergedAt  time.Time
}

// Release is a GitHub GraphQL API Release object.
// See https://docs.github.com/en/graphql/reference/objects#release
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
