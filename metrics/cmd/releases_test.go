package cmd

import (
	"path/filepath"
	"testing"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/test"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
)

func TestReleases(t *testing.T) {
	repo := &github.Repo{Owner: "hackebrot", Name: "turtle"}

	env := map[string]string{
		config.EnvKey("GITHUB", "REPO_OWNER"): "",
		config.EnvKey("GITHUB", "REPO_NAME"):  "",
	}

	tempDir := t.TempDir()

	tests := []test.TestCase{{
		Name:        "releases__repo_owner__required",
		Args:        []string{"github", "-n", repo.Name, "releases"},
		ErrContains: "repo.Owner and repo.Name are required. Set env vars or pass flags",
		Env:         env,
	}, {
		Name:        "releases__repo_name__required",
		Args:        []string{"github", "-o", repo.Owner, "releases"},
		ErrContains: "repo.Owner and repo.Name are required. Set env vars or pass flags",
		Env:         env,
	}, {
		Name:        "releases__default",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "releases"},
		WantFixture: test.NewFixture("github", "releases", "want__default.json"),
		Env:         env,
	}, {
		Name:        "releases__limit",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "releases", "-l", "1"},
		WantFixture: test.NewFixture("github", "releases", "want__limit.json"),
		Env:         env,
	}, {
		Name:        "releases__json",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "releases", "-e", "json"},
		WantFixture: test.NewFixture("github", "releases", "want__default.json"),
		Env:         env,
	}, {
		Name:        "releases__csv",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "releases", "-e", "csv"},
		WantFixture: test.NewFixture("github", "releases", "want__default.csv"),
		Env:         env,
	}, {
		Name:        "releases__filename",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "releases", "-f", filepath.Join(tempDir, "r.json")},
		WantFixture: test.NewFixture("github", "releases", "want__default.json"),
		WantFile:    filepath.Join(tempDir, "r.json"),
		Env:         env,
	}, {
		Name:        "releases__prs__json",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "releases", "--prs", "-e", "json"},
		WantFixture: test.NewFixture("github", "releases", "want__prs.json"),
		Env:         env,
	}, {
		Name:        "releases__prs__csv",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "releases", "--prs", "-e", "csv"},
		WantFixture: test.NewFixture("github", "releases", "want__prs.csv"),
		Env:         env,
	}}

	test.RunTests(t, NewRootCmd, tests)
}
