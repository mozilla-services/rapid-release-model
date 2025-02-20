package github

import (
	"time"
)

// Represents a GitHub repo
type Repo struct {
	Owner string
	Name  string
}

// Unified Git Commit parent model (used for both REST & GraphQL API)
type CommitParent struct {
	AbbreviatedSHA string
	SHA            string
}

// Unified Git Commit Model (used for both REST & GraphQL API)
type Commit struct {
	AbbreviatedSHA string
	SHA            string
	AuthoredDate   time.Time
	CommittedDate  time.Time
	Message        string
	Parents        []*CommitParent
}

// Represents a Git Commit Comparison
type CommitsComparison struct {
	TotalCommits *int
	Commits      []*Commit
}

// Unified GitHub Deployment Model (used for both REST & GraphQL API)
type Deployment struct {
	Description         string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	OriginalEnvironment string
	LatestEnvironment   string
	Task                string
	State               string
	Ref                 string
	Commit              *Commit
}

// Unified GitHub PullRequest Model (used for both REST & GraphQL API)
type PullRequest struct {
	Number    int
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  time.Time
	MergedAt  time.Time
}

// Unified GitHub Release Model (used for both REST & GraphQL API)
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
