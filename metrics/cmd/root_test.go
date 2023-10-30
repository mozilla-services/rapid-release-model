package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	// Name of the test case
	name string

	// Environment variables to be set for the test case
	env map[string]string

	// Arguments to be passed to the CLI app
	args []string

	// Expected output from the CLI app
	output string

	// Text expected in error. Empty string means no error expected.
	errContains string
}

func TestRoot(t *testing.T) {
	repo := &Repo{Owner: "hackebrot", Name: "turtle"}

	t.Setenv("RRM_METRICS_REPO_OWNER", "")
	t.Setenv("RRM_METRICS_REPO_NAME", "")

	tests := []testCase{{
		name:        "repo__short",
		args:        []string{"-o", repo.Owner, "-n", repo.Name},
		output:      fmt.Sprintf("Retrieving metrics for %v/%v", repo.Owner, repo.Name),
		errContains: "",
	}, {
		name:        "repo__long",
		args:        []string{"--repo-owner", repo.Owner, "--repo-name", repo.Name},
		output:      fmt.Sprintf("Retrieving metrics for %v/%v", repo.Owner, repo.Name),
		errContains: "",
	}, {
		name:        "read_env__name",
		env:         map[string]string{"RRM_METRICS_REPO_NAME": "hello"},
		args:        []string{"--repo-owner", repo.Owner},
		output:      fmt.Sprintf("Retrieving metrics for %v/%v", repo.Owner, "hello"),
		errContains: "",
	}, {
		name:        "read_env__owner",
		env:         map[string]string{"RRM_METRICS_REPO_OWNER": "mozilla"},
		args:        []string{"--repo-name", repo.Name},
		output:      fmt.Sprintf("Retrieving metrics for %v/%v", "mozilla", repo.Name),
		errContains: "",
	}, {
		name:        "missing_value__owner",
		args:        []string{"--repo-name", repo.Name},
		errContains: "Repo.Owner is required",
	}, {
		name:        "missing_value__name",
		args:        []string{"--repo-owner", repo.Owner},
		errContains: "Repo.Name is required",
	}}

	runTests(t, tests)
}

// executeCmd creates and executes the metrics rood cmd
func executeCmd(args []string) (string, error) {
	buf := new(bytes.Buffer)

	root := newRootCmd(buf)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err := root.Execute()

	return buf.String(), err
}

// runTests is a helper for table-driven tests using subtests
func runTests(t *testing.T, tests []testCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("running: metrics %s", strings.Join(tt.args, " "))

			if tt.env != nil {
				t.Logf("using environment: %s", tt.env)
				for k, v := range tt.env {
					t.Setenv(k, v)
				}
			}

			got, err := executeCmd(tt.args)

			if tt.errContains != "" && err == nil {
				t.Fatalf("cmd did not return an error. output: %v", got)
			}

			if tt.errContains == "" && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if tt.errContains != "" && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Fatalf("error did not contain message\ngot:     %v\nmissing: %v", err, tt.errContains)
			}

			if tt.output != "" {
				tGot := strings.TrimSpace(got)
				tWant := strings.TrimSpace(tt.output)

				if !cmp.Equal(tGot, tWant) {
					t.Fatalf("cmd returned unexpected output\ngot:  %v\nwant: %v", tGot, tWant)
				}
			}
		})
	}
}
