package cmd

import (
	"testing"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/test"
)

func TestPullRequests(t *testing.T) {
	repo := &github.Repo{Owner: "hackebrot", Name: "turtle"}

	env := map[string]string{
		config.Key("GITHUB", "REPO_OWNER"): "",
		config.Key("GITHUB", "REPO_NAME"):  "",
	}

	tests := []test.TestCase{{
		Name:        "prs__default",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "prs"},
		WantFixture: test.NewFixture("prs", "want__default.json"),
		Env:         env,
	}, {
		Name:        "prs__limit",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "prs", "-l", "2"},
		WantFixture: test.NewFixture("prs", "want__limit.json"),
		Env:         env,
	}, {
		Name:        "prs__csv",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "prs", "-e", "csv"},
		WantFixture: test.NewFixture("prs", "want__default.csv"),
		Env:         env,
	}}

	test.RunTests(t, newRootCmd, tests)
}
