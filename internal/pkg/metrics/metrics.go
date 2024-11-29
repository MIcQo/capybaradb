package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

func StartMetricsServer(port uint) {
	logrus.WithField("port", port).
		Debug("Starting metrics server")

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logrus.Fatal(err)
	}
}
