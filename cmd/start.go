/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"capybaradb/internal/pkg/metrics"
	"capybaradb/internal/pkg/tcp"
	"github.com/sirupsen/logrus"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	// defaultMetricsPort is the default port for the metrics server
	defaultMetricsPort = 8080

	// defaultDatabasePort is the default port for the database server
	defaultDatabasePort = 2121
)

// startServerCmd represents the server command
var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the CapybaraDB server",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting CapybaraDB server...")

		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		var metricsPort, _ = cmd.Flags().GetUint("metricsPort")
		var databasePort, _ = cmd.Flags().GetUint("port")

		go metrics.StartMetricsServer(metricsPort)

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
	startServerCmd.PersistentFlags().Uint("metricsPort", defaultMetricsPort, "Port for the metrics server")
	startServerCmd.PersistentFlags().Uint("port", defaultDatabasePort, "Port for the database server")

}
