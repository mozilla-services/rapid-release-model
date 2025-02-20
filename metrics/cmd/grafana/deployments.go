package grafana

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/grafana"
	"github.com/spf13/cobra"
)

type deploymentsOptions struct {
	App   string
	From  string
	To    string
	Limit int
}

type deploymentsConfig struct {
	*grafanaConfig
	filter *grafana.AnnotationsFilter
}

func newDeploymentsCmd(f Factory, c *grafanaConfig) *cobra.Command {
	opts := new(deploymentsOptions)
	config := &deploymentsConfig{grafanaConfig: c}

	cmd := &cobra.Command{
		Use:   "deployments",
		Short: "Retrieve data about Deployments from Grafana Annotations",
		Long:  "Retrieve data about Deployments from Grafana Annotations",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Read annotations filter values from environment variables.
			// Order: CLI flag default values, environment variables, CLI flag values
			filter := f.DefaultGrafanaAnnotationsFilter()

			// This flag is required. If neither env is set or flag is given, error out.
			if cmd.Flags().Changed("app-name") {
				filter.App = opts.App
			}

			if filter.App == "" {
				return fmt.Errorf("app is required. Set env var or pass flag")
			}

			if filter.From == "" || cmd.Flags().Changed("from") {
				filter.From = opts.From
			}

			if filter.To == "" || cmd.Flags().Changed("to") {
				filter.To = opts.To
			}

			filter.Limit = opts.Limit

			if filter.Limit < 1 {
				return fmt.Errorf("limit cannot be smaller than 1")
			}

			config.filter = filter

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runDeployments(ctx, config)
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.App, "app-name", "a", "", "name of the Grafana app")
	cmd.Flags().StringVar(&opts.From, "from", "now-6M", "epoch datetime in milliseconds (e.g. now-6M)")
	cmd.Flags().StringVar(&opts.To, "to", "now", "epoch datetime in milliseconds (e.g. now)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 100, "limit for how many Deployments to fetch")

	return cmd
}

func runDeployments(ctx context.Context, config *deploymentsConfig) error {
	config.logger.Debug(
		"runDeployments",
		"client", config.client,
	)

	deployments, err := grafana.QueryDeployments(ctx, config.client, config.filter)
	if err != nil {
		return err
	}
	return config.exporter.Export(deployments)
}
