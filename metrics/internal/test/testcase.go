package test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/spf13/cobra"
)

type TestCase struct {
	// Name of the test case
	Name string

	// Environment variables to be set for the test case
	Env map[string]string

	// Arguments to be passed to the CLI app
	Args []string

	// Expected output from the CLI app
	WantText string

	// Function to load extected output
	WantFixture func() ([]byte, error)

	// Text expected in error. Empty string means no error expected.
	ErrContains string
}

// RunTests is a helper function for table-driven tests using subtests
func RunTests(t *testing.T, newCmd func(*factory.Factory) *cobra.Command, tests []TestCase) {
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
			got, err := ExecuteCmd(newCmd, tt.Args)

			if tt.ErrContains != "" && err == nil {
				t.Fatalf("cmd did not return an error. output: %v", got)
			}

			if tt.ErrContains == "" && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if tt.ErrContains != "" && err != nil && !strings.Contains(err.Error(), tt.ErrContains) {
				t.Fatalf("error did not contain message\ngot:     %v\nmissing: %v", err, tt.ErrContains)
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
					t.Fatalf("cmd returned unexpected output\ngot:  %v\nwant: %v", tGot, tWant)
				}
			}
		})
	}
}
