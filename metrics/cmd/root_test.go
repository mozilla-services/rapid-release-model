package cmd

import (
	"testing"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/test"
)

func TestRoot(t *testing.T) {
	env := map[string]string{
		config.EnvKey("GITHUB", "REPO_OWNER"): "",
		config.EnvKey("GITHUB", "REPO_NAME"):  "",
	}

	tests := []test.TestCase{
		{
			Name:        "metrics",
			Args:        []string{},
			ErrContains: "",
			Env:         env,
		},
	}

	test.RunTests(t, newRootCmd, tests)
}
