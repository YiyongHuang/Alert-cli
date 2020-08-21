package cmd

import (
	metricsreport "github.com/YiyongHuang/Alert-cli/pkg/metrics"
	"github.com/spf13/cobra"
)

var cfg metricsreport.MetricCfg

var metricCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Qtt metrics alert",
	Long:  "Qtt metrics alert",
	Run: func(cmd *cobra.Command, args []string) {
		var metricCli metricsreport.MetricsCli
		metricCli.HandleMetrics(&cfg)
	},
}

func init() {
	metricCmd.Flags().StringVar(&cfg.ReportPath, "report-path", "", "The path for send warning")
	metricCmd.Flags().StringVar(&cfg.ReportPathBak, "report-backup-path", "", "The backup path for send warning")
	metricCmd.Flags().StringVar(&cfg.ThanosQueryURL, "thanos-query-url", "", "The path for send warning")
	metricCmd.Flags().StringVar(&cfg.ServicePath, "service-path", "", "The path for get opsservice info")
	rootCmd.AddCommand(metricCmd)
}
