package cmd

import (
	"capybaradb/internal/pkg/config"
	"capybaradb/internal/pkg/metrics"
	"capybaradb/internal/pkg/tcp"
	"capybaradb/internal/pkg/version"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var serverGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "capybaradb",
	Name:      "info",
	Help:      "Shows the CapybaraDB version",
}, []string{"version", "build_date", "code_name", "go_version", "go_os", "go_arch"})

// startServerCmd represents the server command
var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the CapybaraDB server",
	Run: func(cmd *cobra.Command, _ []string) {
		var info = version.AppInfo()
		serverGauge.
			WithLabelValues(
				info.Version,
				info.BuildDate,
				info.Codename,
				info.GoVersion,
				info.GoOS,
				info.GoArch,
			).Add(1)

		logrus.WithField("version", info.Version).
			Info("Starting CapybaraDB server...")

		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		var metricsPort, _ = cmd.Flags().GetUint("metricsPort")
		var metricsEndpoint, _ = cmd.Flags().GetString("metricsEndpoint")
		var databasePort, _ = cmd.Flags().GetUint("port")

		go func() {
			if err := metrics.NewServer(metricsPort, metricsEndpoint).Start(); err != nil {
				logrus.Fatal(err)
			}
		}()

		go func() {
			if err := tcp.NewServer(databasePort).Start(); err != nil {
				logrus.Fatal(err)
			}
		}()

		<-ctx.Done()
		logrus.Println("got interruption signal")
	},
}

func init() {
	rootCmd.AddCommand(startServerCmd)

	// Here you will define your flags and configuration settings.
	startServerCmd.PersistentFlags().Uint("port", config.DefaultDatabasePort, "Port for the database server")

	startServerCmd.PersistentFlags().Uint("metricsPort", config.DefaultMetricsPort, "Port for the metrics server")
	startServerCmd.PersistentFlags().String("metricsEndpoint", config.DefaultMetricsEndpoint, "Endpoint for the metrics server")
}
