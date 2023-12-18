package github

import (
	"regexp"
	"strconv"
)

// ReleaseWithPRs embeds a GitHub GraphQL API Release object with an added slice
// of PR numbers, which we have parsed from its auto-generated release notes.
type ReleaseWithPRs struct {
	Release
	PRs []int
}

// NewReleaseWithPrs creates a new ReleaseWithPRs by parsing PR numbers from
// auto-generated Release Descriptions.
func NewReleaseWithPRs(r *Release) *ReleaseWithPRs {
	var prs []int

	// Pattern for auto-generated release notes. For more information see:
	// https://docs.github.com/en/repositories/releasing-projects-on-github/automatically-generated-release-notes
	re := regexp.MustCompile(`\* .*by @\w+ in .+\/pull\/(?P<pr>\d+)`)

	for _, match := range re.FindAllStringSubmatch(r.Description, -1) {
		for i, name := range re.SubexpNames() {
			if name == "pr" {
				n, err := strconv.Atoi(match[i])
				if err != nil {
					continue
				}
				prs = append(prs, n)
			}
		}
	}

	return &ReleaseWithPRs{*r, prs}
}
