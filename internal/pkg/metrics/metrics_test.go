package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsEndpoint(t *testing.T) {
	// Setup example metrics
	opsProcessed := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "example_processed_ops_total",
		Help: "The total number of processed events",
	})
	prometheus.MustRegister(opsProcessed)

	// Create testing server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create GET ednpoint
	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Could not send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Test if we get 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", resp.StatusCode)
	}

	// Check body and content-type
	expectedContentType := "text/plain; version=0.0.4; charset=utf-8; escaping=values"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("Expected Content-Type %q, got %q", expectedContentType, resp.Header.Get("Content-Type"))
	}
}
