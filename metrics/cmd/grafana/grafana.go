package grafana

import (
	"github.com/mozilla-services/rapid-release-model/metrics/internal/factory"
	"github.com/spf13/cobra"
)

func NewGrafanaCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grafana",
		Short: "Retrieve metrics from Grafana",
		Long:  "Retrieve metrics from Grafana",
	}

	cmd.AddCommand(newDeploymentsCmd(f))

	return cmd
}
