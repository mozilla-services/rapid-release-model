package factory

import (
	"context"
	"reflect"
	"testing"
)

func TestNewDefaultFactory_FunctionFieldsNotNil(t *testing.T) {
	ctx := context.Background()
	f := NewDefaultFactory(ctx) // Directly test the returned instance

	tests := []struct {
		name string
		fn   interface{}
	}{
		{"NewLogger", f.NewLogger},
		{"NewExporter", f.NewExporter},
		{"newEncoder", f.newEncoder},
		{"newGitHubRepo", f.newGitHubRepo},
		{"newGitHubHTTPClient", f.newGitHubHTTPClient},
		{"NewGitHubRESTClient", f.NewGitHubRESTClient},
		{"newGitHubRESTAPI", f.newGitHubRESTAPI},
		{"NewGitHubGraphQLClient", f.NewGitHubGraphQLClient},
		{"newGitHubGraphQLAPI", f.newGitHubGraphQLAPI},
		{"NewGrafanaHTTPClient", f.NewGrafanaHTTPClient},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Factory field: %+v", tt.fn)

			// Use reflection to correctly check for nil function pointers
			if tt.fn == nil || reflect.ValueOf(tt.fn).IsNil() {
				t.Errorf("Factory function %q is nil, expected it to be set", tt.name)
			}
		})
	}
}
