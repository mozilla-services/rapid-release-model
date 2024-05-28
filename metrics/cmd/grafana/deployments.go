package grafana

import (
	"context"
	"fmt"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/grafana"
	"github.com/spf13/cobra"
)

type DeploymentsOptions struct {
	App   string
	From  string
	To    string
	Limit int
}

func newDeploymentsCmd(f *factory.Factory) *cobra.Command {
	opts := new(DeploymentsOptions)

	cmd := &cobra.Command{
		Use:   "deployments",
		Short: "Retrieve data about Deployments from Grafana Annotations",
		Long:  "Retrieve data about Deployments from Grafana Annotations",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Read annotations filter values from environment variables.
			// Order: CLI flag default values, environment variables, CLI flag values
			filter, err := f.NewGrafanaAnnotationsFilter()
			if err != nil {
				return fmt.Errorf("Error reading Grafana Annotations Filter from env: %w", err)
			}

			// This flag is required. If neither env is set or flag is given, error out.
			if cmd.Flags().Changed("app-name") {
				filter.App = opts.App
			}

			if filter.App == "" {
				return fmt.Errorf("App is required. Set env var or pass flag.")
			}

			if filter.From == "" || cmd.Flags().Changed("from") {
				filter.From = opts.From
			}

			if filter.To == "" || cmd.Flags().Changed("to") {
				filter.To = opts.To
			}

			filter.Limit = opts.Limit

			if filter.Limit < 1 {
				return fmt.Errorf("Limit cannot be smaller than 1.")
			}

			// Overwrite the factory function to return the loaded filter values.
			f.NewGrafanaAnnotationsFilter = func() (*grafana.AnnotationsFilter, error) {
				return filter, nil
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployments(cmd.Root().Context(), f, opts)
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.App, "app-name", "a", "", "name of the Grafana app")
	cmd.Flags().StringVar(&opts.From, "from", "now-6M", "epoch datetime in milliseconds (e.g. now-6M)")
	cmd.Flags().StringVar(&opts.To, "to", "now", "epoch datetime in milliseconds (e.g. now)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 100, "limit for how many Deployments to fetch")

	return cmd
}

func runDeployments(ctx context.Context, f *factory.Factory, opts *DeploymentsOptions) error {
	filter, err := f.NewGrafanaAnnotationsFilter()
	if err != nil {
		return err
	}

	client, err := f.NewGrafanaHTTPClient()
	if err != nil {
		return err
	}

	deployments, err := grafana.QueryDeployments(ctx, client, filter)
	if err != nil {
		return err
	}

	exporter, err := f.NewExporter()
	if err != nil {
		return err
	}

	return exporter.Export(deployments)
}
