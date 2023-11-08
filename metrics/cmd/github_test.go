package cmd

import (
	"testing"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/test"
)

func TestGitHub(t *testing.T) {
	env := map[string]string{
		config.Key("GITHUB", "REPO_OWNER"): "",
		config.Key("GITHUB", "REPO_NAME"):  "",
	}

	tests := []test.TestCase{
		{
			Name:        "github",
			Args:        []string{"github"},
			ErrContains: "",
			Env:         env,
		},
	}

	test.RunTests(t, newRootCmd, tests)
}
