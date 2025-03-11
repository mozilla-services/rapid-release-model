package github_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	ghrest "github.com/google/go-github/v68/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github/graphql"
	"github.com/mozilla-services/rapid-release-model/pkg/github/internal/test"
	"github.com/mozilla-services/rapid-release-model/pkg/github/rest"
)

func TestQueryDeployedCommits(t *testing.T) {
	ctx := context.Background()

	logbuf := new(bytes.Buffer)
	logger := slog.New(slog.NewTextHandler(logbuf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	graphQLClient := test.NewFakeGraphQLClient()
	graphQLAPI := graphql.NewGitHubGraphQLAPI(graphQLClient, logger)
	registerGraphQLresponses(t, graphQLClient)

	restClient := test.NewFakeGitHubRESTClient()
	restAPI := rest.NewGitHubRESTAPI(restClient, logger)
	registerRESTresponses(t, restClient)

	tests := []struct {
		name        string
		repo        *github.Repo
		opts        *github.DeployedCommitsOptions
		want        *github.DeploymentWithCommits
		errContains string
	}{
		{
			name: "success",
			repo: &github.Repo{"hackebrot", "turtle"},
			opts: &github.DeployedCommitsOptions{
				Deployment: &github.DeploymentOpts{
					Env:         "stage",
					Sha:         "1abc111aaaaaaaaaaa",
					SearchLimit: 5,
				},
				Commits: &github.CommitsOpts{
					Limit: 250,
				},
			},
			want: &github.DeploymentWithCommits{
				Deployment: &github.Deployment{
					Description:         "Deployment03",
					CreatedAt:           time.Date(2022, time.May, 1, 20, 20, 5, 0, time.UTC),
					UpdatedAt:           time.Date(2022, time.May, 1, 20, 20, 5, 0, time.UTC),
					OriginalEnvironment: "stage",
					LatestEnvironment:   "stage",
					Task:                "deploy",
					State:               "ACTIVE",
					Ref:                 "",
					Commit: &github.Commit{
						AbbreviatedSHA: "1abc111",
						SHA:            "1abc111aaaaaaaaaaa",
						AuthoredDate:   time.Date(2022, time.May, 1, 20, 18, 5, 0, time.UTC),
						CommittedDate:  time.Date(2022, time.May, 1, 20, 18, 5, 0, time.UTC),
						Message:        "commit changes 333",
					},
				},
				DeployedCommits: []*github.Commit{
					{
						AbbreviatedSHA: "1abc111",
						SHA:            "1abc111aaaaaaaaaaa",
						AuthoredDate:   time.Date(2022, time.May, 1, 20, 18, 5, 0, time.UTC),
						CommittedDate:  time.Date(2022, time.May, 1, 20, 18, 5, 0, time.UTC),
						Message:        "commit changes 333",
					},
					{
						AbbreviatedSHA: "2abc111",
						SHA:            "2abc111bbbbbbbbbbb",
						AuthoredDate:   time.Date(2022, time.April, 1, 20, 24, 5, 0, time.UTC),
						CommittedDate:  time.Date(2022, time.April, 1, 20, 24, 5, 0, time.UTC),
						Message:        "commit changes 2222",
					},
				}},
		},
		{
			name: "success__paginated",
			repo: &github.Repo{"hackebrot", "turtle"},
			opts: &github.DeployedCommitsOptions{
				Deployment: &github.DeploymentOpts{
					Env:         "stage",
					Sha:         "2abc111bbbbbbbbbbb",
					SearchLimit: 5,
				},
				Commits: &github.CommitsOpts{
					Limit: 250,
				},
			},
			want: &github.DeploymentWithCommits{
				Deployment: &github.Deployment{
					Description:         "Deployment02",
					CreatedAt:           time.Date(2022, time.April, 1, 20, 25, 5, 0, time.UTC),
					UpdatedAt:           time.Date(2022, time.April, 1, 20, 25, 5, 0, time.UTC),
					OriginalEnvironment: "stage",
					LatestEnvironment:   "stage",
					Task:                "deploy",
					State:               "INACTIVE",
					Ref:                 "",
					Commit: &github.Commit{
						AbbreviatedSHA: "2abc111",
						SHA:            "2abc111bbbbbbbbbbb",
						AuthoredDate:   time.Date(2022, time.April, 1, 20, 24, 5, 0, time.UTC),
						CommittedDate:  time.Date(2022, time.April, 1, 20, 24, 5, 0, time.UTC),
						Message:        "commit changes 2222",
					},
				},
				DeployedCommits: []*github.Commit{
					{
						AbbreviatedSHA: "2abc111",
						SHA:            "2abc111bbbbbbbbbbb",
						AuthoredDate:   time.Date(2022, time.April, 1, 20, 24, 5, 0, time.UTC),
						CommittedDate:  time.Date(2022, time.April, 1, 20, 24, 5, 0, time.UTC),
						Message:        "commit changes 2222",
					},
					{
						AbbreviatedSHA: "5abc111",
						SHA:            "5abc111yyyyyyyyyyy",
						AuthoredDate:   time.Date(2022, time.April, 1, 18, 22, 1, 0, time.UTC),
						CommittedDate:  time.Date(2022, time.April, 1, 18, 22, 1, 0, time.UTC),
						Message:        "commit 3",
					},
					{
						AbbreviatedSHA: "8abc222",
						SHA:            "8abc222eeeeeeeeeee",
						AuthoredDate:   time.Date(2022, time.March, 8, 10, 5, 2, 0, time.UTC),
						CommittedDate:  time.Date(2022, time.March, 8, 10, 5, 2, 0, time.UTC),
						Message:        "commit 2",
					},
					{
						AbbreviatedSHA: "3abc111",
						SHA:            "3abc111ccccccccccc",
						AuthoredDate:   time.Date(2022, time.February, 1, 18, 25, 5, 0, time.UTC),
						CommittedDate:  time.Date(2022, time.February, 1, 18, 25, 5, 0, time.UTC),
						Message:        "commit changes",
					},
				}},
		},
		{
			name: "error__limit",
			repo: &github.Repo{"hackebrot", "turtle"},
			opts: &github.DeployedCommitsOptions{
				Deployment: &github.DeploymentOpts{
					Env:         "stage",
					Sha:         "3abc111ccccccccccc",
					SearchLimit: 2,
				},
				Commits: &github.CommitsOpts{
					Limit: 250,
				},
			},
			errContains: "error querying deployments: search limit 2 reached, no deployment found for SHA 3abc111ccccccccccc in stage",
		},
		{
			name: "error__prev_limit",
			repo: &github.Repo{"hackebrot", "turtle"},
			opts: &github.DeployedCommitsOptions{
				Deployment: &github.DeploymentOpts{
					Env:         "stage",
					Sha:         "3abc111ccccccccccc",
					SearchLimit: 10,
				},
				Commits: &github.CommitsOpts{
					Limit: 250,
				},
			},
			errContains: "error querying deployments: found deployment but no previous deployment for SHA 3abc111ccccccccccc in stage",
		},
		{
			name: "error__prev_deployment",
			repo: &github.Repo{"hackebrot", "turtle"},
			opts: &github.DeployedCommitsOptions{
				Deployment: &github.DeploymentOpts{
					Env:         "prod",
					Sha:         "1abc111aaaaaaaaaaa",
					SearchLimit: 10,
				},
				Commits: &github.CommitsOpts{
					Limit: 250,
				},
			},
			errContains: "found deployment but no previous deployment for SHA 1abc111aaaaaaaaaaa in prod",
		},
		{
			name: "error__nope",
			repo: &github.Repo{"hackebrot", "turtle"},
			opts: &github.DeployedCommitsOptions{
				Deployment: &github.DeploymentOpts{
					Env:         "helloworld",
					Sha:         "1abc111aaaaaaaaaaa",
					SearchLimit: 10,
				},
				Commits: &github.CommitsOpts{
					Limit: 250,
				},
			},
			errContains: "error querying deployments: nope",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// Validate testcase configuration
			if tt.errContains != "" && tt.want != nil {
				t.Fatal("cannot set both errContains and want")
			}

			got, err := github.QueryDeployedCommits(ctx, tt.repo, graphQLAPI, restAPI, logger, tt.opts)

			if tt.errContains != "" && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("error did not contain message\ngot:     %v\nmissing: %v", err, tt.errContains)
				return
			}

			if tt.errContains == "" && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !cmp.Equal(got, tt.want) {
				t.Logf("logging %s\n", logbuf.String())
				t.Errorf("QueryDeployedCommits() = \n  %v\n", cmp.Diff(got, tt.want))
			}

			t.Cleanup(func() {
				logbuf.Reset()
			})
		})
	}
}

func registerGraphQLresponses(t *testing.T, c *test.FakeGraphQLClient) {
	t.Helper()

	deploymentsJsonData1 := `{
		"Repository": {
			"Name": "turtle",
			"Owner": {
				"Login": "hackebrot"
			},
			"Deployments": {
				"PageInfo": {
					"HasNextPage": true,
					"EndCursor": "abc123"
				},
				"Nodes": [
					{
						"Description": "Deployment03",
						"CreatedAt": "2022-05-01T20:20:05Z",
						"UpdatedAt": "2022-05-01T20:20:05Z",
						"OriginalEnvironment": "stage",
						"LatestEnvironment": "stage",
						"Task": "deploy",
						"State": "ACTIVE",
						"Commit": {
							"AbbreviatedOid": "1abc111",
							"Oid": "1abc111aaaaaaaaaaa",
							"AuthoredDate": "2022-05-01T20:18:05Z",
							"CommittedDate": "2022-05-01T20:18:05Z",
							"Message": "commit changes 333"
						}
					},
					{
						"Description": "Deployment02",
						"CreatedAt": "2022-04-01T20:25:05Z",
						"UpdatedAt": "2022-04-01T20:25:05Z",
						"OriginalEnvironment": "stage",
						"LatestEnvironment": "stage",
						"Task": "deploy",
						"State": "INACTIVE",
						"Commit": {
							"AbbreviatedOid": "2abc111",
							"Oid": "2abc111bbbbbbbbbbb",
							"AuthoredDate": "2022-04-01T20:24:05Z",
							"CommittedDate": "2022-04-01T20:24:05Z",
							"Message": "commit changes 2222"
						}
					}
				]
			}
		}
	}`

	c.RegisterResponse(
		test.GraphQLQueryKey{
			QueryType: "*graphql.DeployedCommitsQuery",
			RepoOwner: "hackebrot",
			RepoName:  "turtle",
			Extra:     test.GraphQLQueryKeyExtra{Environments: "stage"},
			EndCursor: "",
		},
		&test.GraphQLResponse{
			Content: deploymentsJsonData1,
		},
	)

	deploymentsJsonData2 := `{
	    "Repository": {
			"Name": "turtle",
			"Owner": {
				"Login": "hackebrot"
			},
			"Deployments": {
				"PageInfo": {
					"HasNextPage": false,
					"EndCursor": "abc456"
				},
				"Nodes": [
					{
						"Description": "Deployment01",
						"CreatedAt": "2022-02-01T20:25:05Z",
						"UpdatedAt": "2022-02-01T20:25:05Z",
						"OriginalEnvironment": "stage",
						"LatestEnvironment": "stage",
						"Task": "deploy",
						"State": "INACTIVE",
						"Commit": {
							"AbbreviatedOid": "3abc111",
							"Oid": "3abc111ccccccccccc",
							"AuthoredDate": "2022-02-01T18:25:05Z",
							"CommittedDate": "2022-02-01T18:25:05Z",
							"Message": "commit changes"
						}
					}
				]
			}
		}
	}`

	c.RegisterResponse(
		test.GraphQLQueryKey{
			QueryType: "*graphql.DeployedCommitsQuery",
			RepoOwner: "hackebrot",
			RepoName:  "turtle",
			Extra:     test.GraphQLQueryKeyExtra{Environments: "stage"},
			EndCursor: "abc123",
		},
		&test.GraphQLResponse{
			Content: deploymentsJsonData2,
		},
	)

	deploymentsJsonData3 := `{
	    "Repository": {
			"Name": "turtle",
			"Owner": {
				"Login": "hackebrot"
			},
			"Deployments": {
				"PageInfo": {
					"HasNextPage": false,
					"EndCursor": "abc"
				},
				"Nodes": [
				{
						"Description": "Deployment03",
						"CreatedAt": "2022-05-02T20:25:05Z",
						"UpdatedAt": "2022-05-02T20:25:05Z",
						"OriginalEnvironment": "prod",
						"LatestEnvironment": "prod",
						"Task": "deploy",
						"State": "ACTIVE",
						"Commit": {
							"AbbreviatedOid": "1abc111",
							"Oid": "1abc111aaaaaaaaaaa",
							"AuthoredDate": "2022-05-01T20:18:05Z",
							"CommittedDate": "2022-05-01T20:18:05Z",
							"Message": "commit changes 333"
						}
					}
				]
			}
		}
	}`

	c.RegisterResponse(
		test.GraphQLQueryKey{
			QueryType: "*graphql.DeployedCommitsQuery",
			RepoOwner: "hackebrot",
			RepoName:  "turtle",
			Extra:     test.GraphQLQueryKeyExtra{Environments: "prod"},
			EndCursor: "",
		},
		&test.GraphQLResponse{
			Content: deploymentsJsonData3,
		},
	)

	c.RegisterResponse(
		test.GraphQLQueryKey{
			QueryType: "*graphql.DeployedCommitsQuery",
			RepoOwner: "hackebrot",
			RepoName:  "turtle",
			Extra:     test.GraphQLQueryKeyExtra{Environments: "helloworld"},
			EndCursor: "",
		},
		&test.GraphQLResponse{
			Err: errors.New("nope"),
		},
	)
}

func registerRESTresponses(t *testing.T, c *test.FakeGitHubRESTClient) {
	t.Helper()
	c.RegisterCommitComparison(
		test.RESTCommitComparisonQueryKey{
			RepoOwner: "hackebrot",
			RepoName:  "turtle",
			Base:      "2abc111bbbbbbbbbbb",
			Head:      "1abc111aaaaaaaaaaa",
			Page:      1,
		},
		&test.RESTCommitComparisonResponse{
			APIResponse: &ghrest.Response{NextPage: 0},
			Comparison: &ghrest.CommitsComparison{
				TotalCommits: ghrest.Ptr(2),
				Commits: []*ghrest.RepositoryCommit{
					{
						SHA: ghrest.Ptr("1abc111aaaaaaaaaaa"),
						Commit: &ghrest.Commit{
							SHA:       ghrest.Ptr("1abc111aaaaaaaaaaa"),
							Author:    &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.May, 1, 20, 18, 5, 0, time.UTC)}},
							Committer: &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.May, 1, 20, 18, 5, 0, time.UTC)}},
							Message:   ghrest.Ptr("commit changes 333"),
						},
					},
					{
						SHA: ghrest.Ptr("2abc111bbbbbbbbbbb"),
						Commit: &ghrest.Commit{
							SHA:       ghrest.Ptr("2abc111bbbbbbbbbbb"),
							Author:    &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.April, 1, 20, 24, 5, 0, time.UTC)}},
							Committer: &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.April, 1, 20, 24, 5, 0, time.UTC)}},
							Message:   ghrest.Ptr("commit changes 2222"),
						},
					},
				},
			},
		},
	)

	c.RegisterCommitComparison(
		test.RESTCommitComparisonQueryKey{
			RepoOwner: "hackebrot",
			RepoName:  "turtle",
			Base:      "3abc111ccccccccccc",
			Head:      "2abc111bbbbbbbbbbb",
			Page:      1,
		},
		&test.RESTCommitComparisonResponse{
			APIResponse: &ghrest.Response{NextPage: 2},
			Comparison: &ghrest.CommitsComparison{
				TotalCommits: ghrest.Ptr(3),
				Commits: []*ghrest.RepositoryCommit{
					{
						SHA: ghrest.Ptr("2abc111bbbbbbbbbbb"),
						Commit: &ghrest.Commit{
							SHA:       ghrest.Ptr("2abc111bbbbbbbbbbb"),
							Author:    &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.April, 1, 20, 24, 5, 0, time.UTC)}},
							Committer: &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.April, 1, 20, 24, 5, 0, time.UTC)}},
							Message:   ghrest.Ptr("commit changes 2222"),
						},
					},
					{
						SHA: ghrest.Ptr("5abc111yyyyyyyyyyy"),
						Commit: &ghrest.Commit{
							SHA:       ghrest.Ptr("5abc111yyyyyyyyyyy"),
							Author:    &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.April, 1, 18, 22, 1, 0, time.UTC)}},
							Committer: &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.April, 1, 18, 22, 1, 0, time.UTC)}},
							Message:   ghrest.Ptr("commit 3"),
						},
					},
					{
						SHA: ghrest.Ptr("8abc222eeeeeeeeeee"),
						Commit: &ghrest.Commit{
							SHA:       ghrest.Ptr("8abc222eeeeeeeeeee"),
							Author:    &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.March, 8, 10, 5, 2, 0, time.UTC)}},
							Committer: &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.March, 8, 10, 5, 2, 0, time.UTC)}},
							Message:   ghrest.Ptr("commit 2"),
						},
					},
				},
			},
		},
	)

	c.RegisterCommitComparison(
		test.RESTCommitComparisonQueryKey{
			RepoOwner: "hackebrot",
			RepoName:  "turtle",
			Base:      "3abc111ccccccccccc",
			Head:      "2abc111bbbbbbbbbbb",
			Page:      2,
		},
		&test.RESTCommitComparisonResponse{
			APIResponse: &ghrest.Response{NextPage: 0},
			Comparison: &ghrest.CommitsComparison{
				TotalCommits: ghrest.Ptr(1),
				Commits: []*ghrest.RepositoryCommit{
					{
						SHA: ghrest.Ptr("3abc111ccccccccccc"),
						Commit: &ghrest.Commit{
							SHA:       ghrest.Ptr("3abc111ccccccccccc"),
							Author:    &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.February, 1, 18, 25, 5, 0, time.UTC)}},
							Committer: &ghrest.CommitAuthor{Date: &ghrest.Timestamp{Time: time.Date(2022, time.February, 1, 18, 25, 5, 0, time.UTC)}},
							Message:   ghrest.Ptr("commit changes"),
						},
					},
				},
			},
		},
	)

}
