package cmd

import (
	"testing"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/test"
)

func TestReleases(t *testing.T) {
	repo := &github.Repo{Owner: "hackebrot", Name: "turtle"}

	env := map[string]string{
		config.Key("GITHUB", "REPO_OWNER"): "",
		config.Key("GITHUB", "REPO_NAME"):  "",
	}

	tests := []test.TestCase{{
		Name:        "releases__default",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "releases"},
		WantFixture: test.NewFixture("releases", "want__default.json"),
		Env:         env,
	}, {
		Name:        "releases__limit",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "releases", "-l", "1"},
		WantFixture: test.NewFixture("releases", "want__limit.json"),
		Env:         env,
	}}

	test.RunTests(t, newRootCmd, tests)
}
