package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func NewServer(port uint, endpoint string) *Server {
	return &Server{
		port:     port,
		endpoint: endpoint,
	}
}

type Server struct {
	port     uint
	endpoint string
}

func (s *Server) Start() error {
	logrus.WithField("port", s.port).
		Debug("Starting metrics server")

	http.Handle(s.endpoint, promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
