package cmd

import (
	"fmt"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
	"github.com/spf13/cobra"
)

func newGitHubCmd(f *factory.Factory) *cobra.Command {
	repo, err := f.NewGitHubRepo()
	if err != nil {
		panic(err)
	}

	cmd := &cobra.Command{
		Use:   "github",
		Short: "Retrieve metrics from GitHub",
		Long:  "Retrieve metrics from GitHub",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if repo.Owner == "" || repo.Name == "" {
				return fmt.Errorf("Repo.Owner and Repo.Name are required. Set env vars or pass flags.")
			}
			f.NewGitHubRepo = func() (*github.Repo, error) {
				return repo, nil
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&repo.Owner, "repo-owner", "o", repo.Owner, "owner of the GitHub repo")
	cmd.PersistentFlags().StringVarP(&repo.Name, "repo-name", "n", repo.Name, "name of the GitHub repo")

	cmd.AddCommand(newPullRequestsCmd(f))

	return cmd
}
