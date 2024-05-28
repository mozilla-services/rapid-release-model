package grafana

import (
	"net/url"
	"testing"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/config"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/test"
)

func TestDeployments(t *testing.T) {
	env := map[string]string{
		config.EnvKey("GRAFANA", "TOKEN"):               "",
		config.EnvKey("GRAFANA", "SERVER_URL"):          "",
		config.EnvKey("GRAFANA", "ANNOTATIONS", "APP"):  "",
		config.EnvKey("GRAFANA", "ANNOTATIONS", "FROM"): "",
		config.EnvKey("GRAFANA", "ANNOTATIONS", "TO"):   "",
	}

	tests := []test.TestCase{
		{
			Name:        "deployments__app_name__required",
			Args:        []string{"deployments"},
			ErrContains: "App is required. Set env var or pass flag",
			Env:         env,
		},
		{
			Name:        "deployments__defaults",
			Args:        []string{"deployments", "-a", "turtle"},
			WantFixture: test.NewFixture("api", "annotations", "want__defaults.json"),
			WantReqParams: &test.WantReqParams{Grafana: &test.GrafanaReqParams{
				Path: "api/annotations",
				Params: url.Values{
					"type":  []string{"annotation"},
					"from":  []string{"now-6M"},
					"to":    []string{"now"},
					"limit": []string{"100"},
					"tags":  []string{"event_type:deployment", "event_status:complete", "app:turtle"},
				},
			}},
			Env: env,
		},
		{
			Name: "deployments__env__from_to",
			Args: []string{"deployments", "-a", "turtle"},
			WantReqParams: &test.WantReqParams{Grafana: &test.GrafanaReqParams{
				Path: "api/annotations",
				Params: url.Values{
					"type":  []string{"annotation"},
					"from":  []string{"now-4M"},
					"to":    []string{"now-2M"},
					"limit": []string{"100"},
					"tags":  []string{"event_type:deployment", "event_status:complete", "app:turtle"},
				},
			}},
			Env: map[string]string{
				config.EnvKey("GRAFANA", "ANNOTATIONS", "FROM"): "now-4M",
				config.EnvKey("GRAFANA", "ANNOTATIONS", "TO"):   "now-2M",
			},
		},
		{
			Name: "deployments__env__from_to__overwrite",
			Args: []string{"deployments", "-a", "turtle", "--from", "now-12M", "--to", "now-1M"},
			WantReqParams: &test.WantReqParams{Grafana: &test.GrafanaReqParams{
				Path: "api/annotations",
				Params: url.Values{
					"type":  []string{"annotation"},
					"from":  []string{"now-12M"},
					"to":    []string{"now-1M"},
					"limit": []string{"100"},
					"tags":  []string{"event_type:deployment", "event_status:complete", "app:turtle"},
				},
			}},
			Env: map[string]string{
				config.EnvKey("GRAFANA", "ANNOTATIONS", "FROM"): "now-4M",
				config.EnvKey("GRAFANA", "ANNOTATIONS", "TO"):   "now-2M",
			},
		},
	}

	test.RunTests(t, NewGrafanaCmd, tests)
}
