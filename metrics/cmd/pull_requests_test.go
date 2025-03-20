package cmd

import (
	"path/filepath"
	"testing"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/test"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
)

func TestPullRequests(t *testing.T) {
	repo := &github.Repo{Owner: "hackebrot", Name: "turtle"}

	env := map[string]string{
		config.EnvKey("GITHUB", "REPO_OWNER"): "",
		config.EnvKey("GITHUB", "REPO_NAME"):  "",
	}

	tempDir := t.TempDir()

	tests := []test.TestCase{{
		Name:        "prs__repo_owner__required",
		Args:        []string{"github", "-n", repo.Name, "prs"},
		ErrContains: "repo.Owner and repo.Name are required. Set env vars or pass flags",
		Env:         env,
	}, {
		Name:        "prs__repo_name__required",
		Args:        []string{"github", "-o", repo.Owner, "prs"},
		ErrContains: "repo.Owner and repo.Name are required. Set env vars or pass flags",
		Env:         env,
	}, {
		Name:        "prs__default",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "prs"},
		WantFixture: test.NewFixture("github", "prs", "want__default.json"),
		Env:         env,
	}, {
		Name:        "prs__limit",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "prs", "-l", "2"},
		WantFixture: test.NewFixture("github", "prs", "want__limit.json"),
		Env:         env,
	}, {
		Name:        "prs__json",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "prs", "-e", "json"},
		WantFixture: test.NewFixture("github", "prs", "want__default.json"),
		Env:         env,
	}, {
		Name:        "prs__csv",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "prs", "-e", "csv"},
		WantFixture: test.NewFixture("github", "prs", "want__default.csv"),
		Env:         env,
	}, {
		Name:        "prs__filename",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "prs", "-f", filepath.Join(tempDir, "prs.json")},
		WantFixture: test.NewFixture("github", "prs", "want__default.json"),
		WantFile:    filepath.Join(tempDir, "prs.json"),
		Env:         env,
	}}

	test.RunTests(t, NewRootCmd, tests)
}
