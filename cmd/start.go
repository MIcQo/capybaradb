/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"capybaradb/internal/pkg/metrics"
	"capybaradb/internal/pkg/tcp"
	"capybaradb/internal/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	// defaultDatabasePort is the default port for the database server
	defaultDatabasePort = 2121
)

const (
	// defaultMetricsEndpoint is the default endpoint for the metrics server
	defaultMetricsEndpoint = "/metrics"

	// defaultMetricsPort is the default port for the metrics server
	defaultMetricsPort = 8080
)

var ServerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "capybaradb",
	Name:      "info",
	Help:      "Shows the CapybaraDB version",
}, []string{"version", "build_date", "code_name", "go_version", "go_os", "go_arch"})

// startServerCmd represents the server command
var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the CapybaraDB server",
	Run: func(cmd *cobra.Command, args []string) {
		var info = version.AppInfo()
		ServerGauge.
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
	startServerCmd.PersistentFlags().Uint("port", defaultDatabasePort, "Port for the database server")

	startServerCmd.PersistentFlags().Uint("metricsPort", defaultMetricsPort, "Port for the metrics server")
	startServerCmd.PersistentFlags().String("metricsEndpoint", defaultMetricsEndpoint, "Endpoint for the metrics server")
}
