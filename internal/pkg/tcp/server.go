// Package tcp contains the TCP Server
package tcp

import (
	"bufio"
	"bytes"
	"capybaradb/internal/pkg/engine"
	"capybaradb/internal/pkg/user"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
	"vitess.io/vitess/go/vt/sqlparser"
)

var openConnectionsGauge = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace:   "capybaradb",
		Subsystem:   "",
		Name:        "open_connections",
		Help:        "Number of open connections",
		ConstLabels: nil,
	},
	[]string{},
)

// Config is the configuration for the TCP Server
type Config func(*Server)

// WithPort sets the port for the TCP Server
func WithPort(port uint) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

// WithEngineConfig sets the engine configuration
func WithEngineConfig(engineConfig *engine.Config) func(*Server) {
	return func(s *Server) {
		s.engineConfig = engineConfig
	}
}

// NewServer creates a new TCP Server
func NewServer(cfg ...Config) *Server {
	var s = &Server{}

	for _, c := range cfg {
		c(s)
	}

	return s
}

// Server represents the TCP Server
type Server struct {
	port   uint
	server net.Listener

	engineConfig *engine.Config
}

// Start starts the TCP Server
func (s *Server) Start() (err error) {
	logrus.WithField("port", s.port).
		Debug("Starting database engine")

	s.server, err = net.Listen("tcp", s.addr())
	if err != nil {
		return
	}
	return s.handleConnections()
}

// Close shuts down the TCP Server
func (s *Server) Close() (err error) {
	return s.server.Close()
}

func (s *Server) addr() string {
	return fmt.Sprintf(":%d", s.port)
}

func (s *Server) handleConnections() (err error) {
	for {
		conn, acceptErr := s.server.Accept()
		if acceptErr != nil || conn == nil {
			err = errors.New("could not accept connection")
			break
		}

		openConnectionsGauge.WithLabelValues().Inc()
		logrus.WithField("addr", conn.RemoteAddr().String()).Debug("new connection")

		go s.handleConnection(conn)
	}
	return
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	var uctx = &user.Context{Schema: "public"}
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		req, err := rw.ReadString('\n')
		if err != nil {
			_, _ = rw.WriteString("failed to read input")
			_ = rw.Flush()
			break
		}

		queryBytes, err := hex.DecodeString(req[0 : len(req)-1])
		if err != nil {
			logrus.WithError(err).WithField("input", req).Debug("Failed to decode input")
			_, _ = rw.WriteString("failed to decode input")
			_ = rw.Flush()
			continue
		}

		var query = strings.TrimSuffix(strings.TrimSuffix(string(queryBytes), "\n"), "\n")

		var logger = logrus.WithField("query", query)
		logger.Debug("received query")

		parser, err := sqlparser.New(sqlparser.Options{})
		if err != nil {
			logger.WithError(err).Debug("failed to create parser")
			_, _ = rw.WriteString("failed to parse input")
			_ = rw.Flush()
			continue
		}

		stmt, err := parser.Parse(query)
		if err != nil {
			logger.WithError(err).Debug("failed to parse input")
			_, _ = rw.WriteString("failed to parse input")
			_ = rw.Flush()
			continue
		}

		uctx.Query = query
		result, err := engine.ExecuteStatement(s.engineConfig, uctx, stmt)
		//if err != nil {
		//	logger.WithError(err).Debug("failed to execute statement")
		//	_, _ = rw.WriteString("failed to execute statement: " + err.Error())
		//	_ = rw.Flush()
		//	continue
		//}

		var response = new(bytes.Buffer)
		var encoder = gob.NewEncoder(response)
		var data = Packet{
			Error: err != nil,
		}

		if err == nil {
			data.AffectedRows = int64(result.AffectedRows())
			data.LastInsertID = int32(result.LastInsertId())
			data.Rows = result.Rows()
			data.Columns = result.Columns()
		}

		if encodeErr := encoder.Encode(data); encodeErr != nil {
			logger.WithError(encodeErr).Debug("failed to encode response")
			_, _ = rw.WriteString("failed to encode response")
			_ = rw.Flush()
			continue
		}

		_, _ = rw.Write(response.Bytes())

		//logrus.Debugf("response: %+#v", data)
		//
		//binaryWriteErr := binary.Write(response, binary.BigEndian, data)
		//if binaryWriteErr != nil {
		//	logger.WithError(binaryWriteErr).Debug("failed to write response")
		//	_, _ = rw.WriteString("failed to write response")
		//	_ = rw.Flush()
		//	continue
		//}
		//
		//_, _ = rw.Write(response.Bytes())

		//if len(result.Rows()) > 0 {
		//	var t = table.NewWriter()
		//	t.SetOutputMirror(rw)
		//	var headers = table.Row{}
		//	for _, header := range result.Columns() {
		//		headers = append(headers, header)
		//	}
		//
		//	t.AppendHeader(headers)
		//
		//	for _, row := range result.Rows() {
		//		var r = table.Row{}
		//		for _, cell := range row {
		//			r = append(r, cell)
		//		}
		//
		//		t.AppendRow(r)
		//	}
		//
		//	t.Render()
		//} else {
		//	_, _ = rw.WriteString(fmt.Sprintf(
		//		"Affected rows: %d",
		//		result.AffectedRows(),
		//	))
		//}

		_ = rw.Flush()
	}

	logrus.WithField("addr", conn.RemoteAddr().String()).Debug("connection closed")
	openConnectionsGauge.WithLabelValues().Dec()
}

type Packet struct {
	Error        bool
	AffectedRows int64
	LastInsertID int32
	Columns      []string
	Rows         [][]string
}
