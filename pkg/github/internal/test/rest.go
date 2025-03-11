package test

import (
	"context"
	"fmt"
	"sync"

	ghrest "github.com/google/go-github/v68/github"
	"github.com/mozilla-services/rapid-release-model/pkg/github/rest"
)

// RESTCommitComparisonQueryKey uniquely identifies a CompareCommits request.
type RESTCommitComparisonQueryKey struct {
	RepoOwner string
	RepoName  string
	Base      string
	Head      string
	Page      int
}

// CommitComparisonResponse stores the API response for CompareCommits.
type RESTCommitComparisonResponse struct {
	Comparison  *ghrest.CommitsComparison
	APIResponse *ghrest.Response
	Err         error
}

// FakeGitHubRESTClient is a mock REST API client that stores responses in-memory.
type FakeGitHubRESTClient struct {
	mu                 sync.Mutex
	commitsComparisons map[RESTCommitComparisonQueryKey]*RESTCommitComparisonResponse
}

// NewFakeGitHubRESTClient initializes an in-memory fake REST client.
func NewFakeGitHubRESTClient() *FakeGitHubRESTClient {
	return &FakeGitHubRESTClient{
		commitsComparisons: make(map[RESTCommitComparisonQueryKey]*RESTCommitComparisonResponse),
	}
}

// RegisterCommitComparison registers a commit comparison response.
func (f *FakeGitHubRESTClient) RegisterCommitComparison(queryKey RESTCommitComparisonQueryKey, response *RESTCommitComparisonResponse) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.commitsComparisons[queryKey] = response
}

// CompareCommits returns the registered response or an error.
func (f *FakeGitHubRESTClient) CompareCommits(ctx context.Context, owner, repo, base, head string, opts *ghrest.ListOptions) (*ghrest.CommitsComparison, *ghrest.Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	queryKey := RESTCommitComparisonQueryKey{
		RepoOwner: owner,
		RepoName:  repo,
		Base:      base,
		Head:      head,
		Page:      opts.Page,
	}

	// Return registered response
	if response, exists := f.commitsComparisons[queryKey]; exists {
		if response.Err != nil {
			return nil, response.APIResponse, response.Err
		}
		return response.Comparison, response.APIResponse, nil
	}

	return nil, nil, fmt.Errorf("no response registered for query: %v", queryKey)
}

var (
	_ rest.Client = (*FakeGitHubRESTClient)(nil)
)
