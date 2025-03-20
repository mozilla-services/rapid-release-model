package grafana

import (
	"fmt"
	"log/slog"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/export"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/mozilla-services/rapid-release-model/metrics/internal/grafana"
	"github.com/spf13/cobra"
)

type Factory interface {
	factory.GenericFactory
	factory.GrafanaFactory
}

type grafanaConfig struct {
	logger   *slog.Logger
	exporter export.Exporter
	client   grafana.HTTPClient
}

func NewGrafanaCmd(f Factory) *cobra.Command {
	config := new(grafanaConfig)

	cmd := &cobra.Command{
		Use:   "grafana",
		Short: "Retrieve metrics from Grafana",
		Long:  "Retrieve metrics from Grafana",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger, err := f.Logger()
			if err != nil {
				return fmt.Errorf("error retrieving logger: %w", err)
			}
			config.logger = logger

			exporter, err := f.Exporter()
			if err != nil {
				return fmt.Errorf("error retrieving exporter: %w", err)
			}
			config.exporter = exporter

			if err := f.ConfigureGrafanaHubHTTPClient(); err != nil {
				return fmt.Errorf("error configuring Grafana HTTP client: %w", err)
			}

			client, err := f.GrafanaHubHTTPClient()
			if err != nil {
				return fmt.Errorf("error retrieving Grafana HTTP client: %w", err)
			}
			config.client = client

			return nil
		},
	}

	cmd.AddCommand(newDeploymentsCmd(f, config))

	return cmd
}
