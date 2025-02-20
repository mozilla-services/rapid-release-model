package github

import (
	"time"

	"github.com/shurcooL/githubv4"
)

type Commit struct {
	AbbreviatedOid githubv4.GitObjectID
	Oid            githubv4.GitObjectID
	AuthoredDate   time.Time
	CommittedDate  time.Time
	Message        string
}
