package cmd

import (
	"path/filepath"
	"testing"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/test"
	"github.com/shurcooL/githubv4"
)

func TestDeployments(t *testing.T) {
	repo := &github.Repo{Owner: "hackebrot", Name: "turtle"}

	env := map[string]string{
		config.EnvKey("GITHUB", "REPO_OWNER"): "",
		config.EnvKey("GITHUB", "REPO_NAME"):  "",
	}

	tempDir := t.TempDir()

	tests := []test.TestCase{{
		Name:        "deployments__repo_owner__required",
		Args:        []string{"github", "-n", repo.Name, "deployments"},
		ErrContains: "Repo.Owner and Repo.Name are required. Set env vars or pass flags",
		Env:         env,
	}, {
		Name:        "deployments__repo_name__required",
		Args:        []string{"github", "-o", repo.Owner, "deployments"},
		ErrContains: "Repo.Owner and Repo.Name are required. Set env vars or pass flags",
		Env:         env,
	}, {
		Name:        "deployments__default",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "deployments"},
		WantFixture: test.NewFixture("deployments", "want__default.json"),
		Env:         env,
	}, {
		Name:        "deployments__limit",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "deployments", "-l", "2"},
		WantFixture: test.NewFixture("deployments", "want__limit.json"),
		Env:         env,
	}, {
		Name:        "deployments__json",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "deployments", "-e", "json"},
		WantFixture: test.NewFixture("deployments", "want__default.json"),
		Env:         env,
	}, {
		Name:        "deployments__csv",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "deployments", "-e", "csv"},
		WantFixture: test.NewFixture("deployments", "want__default.csv"),
		Env:         env,
	}, {
		Name:        "deployments__filename",
		Args:        []string{"github", "-o", repo.Owner, "-n", repo.Name, "deployments", "-f", filepath.Join(tempDir, "r.json")},
		WantFixture: test.NewFixture("deployments", "want__default.json"),
		WantFile:    filepath.Join(tempDir, "r.json"),
		Env:         env,
	}, {
		Name: "deployments__env__single",
		Args: []string{"github", "-o", repo.Owner, "-n", repo.Name, "deployments", "--env", "prod"},
		WantReqParams: &test.WantReqParams{
			GitHub: &test.GitHubReqParams{
				Variables: map[string]interface{}{
					"environments": []githubv4.String{githubv4.String("prod")}},
			},
		},
		Env: env,
	}, {
		Name: "deployments__env__multiple",
		Args: []string{"github", "-o", repo.Owner, "-n", repo.Name, "deployments", "--env", "prod", "--env", "hello"},
		WantReqParams: &test.WantReqParams{
			GitHub: &test.GitHubReqParams{
				Variables: map[string]interface{}{
					"environments": []githubv4.String{githubv4.String("prod"), githubv4.String("hello")}},
			},
		},
		Env: env,
	}}

	test.RunTests(t, newRootCmd, tests)
}
