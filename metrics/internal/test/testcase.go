package test

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/spf13/cobra"
)

// WantReqParams holds expected values for outgoing requests to GitHub and
// Grafana. These will be passed on to the respective fake clients and will be
// validated against when performing a request.
type WantReqParams struct {
	GitHub  *GitHubReqParams
	Grafana *GrafanaReqParams
}

type TestCase struct {
	// Name of the test case
	Name string

	// Environment variables to be set for the test case
	Env map[string]string

	// Arguments to be passed to the CLI app
	Args []string

	// Expected output from the CLI app
	WantText string

	// Expected log output from the CLI app
	WantLog string

	// Function to load extected output
	WantFixture func() ([]byte, error)

	// Expect output to be written to this file
	WantFile string

	// Expected parameters for outgoing requests to GitHub and Grafana
	WantReqParams *WantReqParams

	// Text expected in error. Empty string means no error expected.
	ErrContains string
}

// RunTests is a helper function for table-driven tests using subtests
func RunTests(t *testing.T, newCmd func(*factory.DefaultFactory) *cobra.Command, tests []TestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Logf("running: metrics %s", strings.Join(tt.Args, " "))

			// Set environment variables
			if tt.Env != nil {
				t.Logf("using environment: %s", tt.Env)
				for k, v := range tt.Env {
					t.Setenv(k, v)
				}
			}

			// Validate testcase configuration
			if (tt.ErrContains != "" && tt.WantFixture != nil) || (tt.ErrContains != "" && tt.WantText != "") {
				t.Fatalf("cannot set both errContains and wantFixture or wantText")
			}

			if tt.WantFixture != nil && tt.WantText != "" {
				t.Fatalf("cannot set both wantFixture and wantText")
			}

			// Execute the CLI cmd with the specified args
			got, log, err := ExecuteCmd(newCmd, tt.Args, tt.WantReqParams)

			if tt.ErrContains != "" && err == nil {
				t.Fatalf("cmd did not return an error. output: %v", got)
			}

			if tt.ErrContains == "" && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if tt.ErrContains != "" && err != nil && !strings.Contains(err.Error(), tt.ErrContains) {
				t.Fatalf("error did not contain message\ngot:     %v\nmissing: %v", err, tt.ErrContains)
			}

			if tt.WantLog != "" {
				logGot := strings.TrimSpace(log)
				logWant := strings.TrimSpace(tt.WantLog)

				if !strings.Contains(logGot, logWant) {
					t.Fatalf("cmd log did not contain expected message\n%v", cmp.Diff(logGot, logWant))
				}
			}

			if tt.WantFile != "" {
				export, err := os.ReadFile(tt.WantFile)
				if err != nil && got == "" {
					t.Fatalf("CLI did not create file at %v. buffer is empty.", tt.WantFile)
				}
				if err != nil && got != "" {
					t.Fatalf("CLI did not create file at %v. buffer:\n%v", tt.WantFile, got)
				}
				// Overwrite got with the contents of the exported file
				got = string(export[:])
			}

			want := tt.WantText

			if tt.WantFixture != nil {
				fixtureData, err := tt.WantFixture()
				if err != nil {
					t.Fatalf("error loading fixture: %v", err)
				}
				want = string(fixtureData[:])
			}

			if want != "" {
				tGot := strings.TrimSpace(got)
				tWant := strings.TrimSpace(want)

				if !cmp.Equal(tGot, tWant) {
					t.Fatalf("cmd returned unexpected output\n%v", cmp.Diff(tGot, tWant))
				}
			}
		})
	}
}
