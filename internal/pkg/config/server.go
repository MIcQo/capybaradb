// Package config hold default values for the server
package config

const (
	// DefaultHost is the default host of the database
	DefaultHost = "127.0.0.1"

	// DefaultDatabasePort is the default port for the database server
	DefaultDatabasePort = 2121

	// DefaultSchema is the default schema
	DefaultSchema = "public"

	// DefaultInputBufferSize defines size of input buffer for TCP packets
	DefaultInputBufferSize = 2048
)

const (
	// DefaultMetricsEndpoint is the default endpoint for the metrics server
	DefaultMetricsEndpoint = "/metrics"

	// DefaultMetricsPort is the default port for the metrics server
	DefaultMetricsPort = 8080
)
