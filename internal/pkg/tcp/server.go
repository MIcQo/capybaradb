// Package tcp contains the TCP Server
package tcp

import (
	"bufio"
	"capybaradb/internal/pkg/config"
	"capybaradb/internal/pkg/engine"
	"capybaradb/internal/pkg/mysql-protocol"
	"capybaradb/internal/pkg/session"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"net"
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

// WithReadBufferSize sets read buffer size for TCP packets
func WithReadBufferSize(bufferSize uint) func(*Server) {
	return func(c *Server) {
		c.readBufferSize = bufferSize
	}
}

// NewServer creates a new TCP Server
func NewServer(cfg ...Config) *Server {
	var s = &Server{}

	for _, c := range cfg {
		c(s)
	}

	if s.readBufferSize == 0 {
		s.readBufferSize = config.DefaultInputBufferSize
	}

	return s
}

// Server represents the TCP Server
type Server struct {
	port           uint
	server         net.Listener
	readBufferSize uint

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
		logrus.WithField("remote", conn.RemoteAddr()).Debug("closing connection")
		openConnectionsGauge.WithLabelValues().Dec()
		_ = conn.Close()
	}(conn)

	var connSession = session.NewContext("")

	logrus.WithField("remote", conn.RemoteAddr()).Debug("new client connected")

	var handshake = mysql_protocol.NewDefaultHandshakePacket()
	var packet = handshake.Encode()
	_, _ = conn.Write(packet)

	// read login packet
	buffer := make([]byte, 4096)
	_, err := conn.Read(buffer)
	if err != nil {
		logrus.Fatal(err)
		return
	}

	var lp, _ = mysql_protocol.NewLoginPacket().Decode(buffer)
	var loginPacket = lp.(*mysql_protocol.LoginPacket)

	var isPasswordValid = mysql_protocol.ValidatePassword("aa", handshake, loginPacket)

	// validate
	if loginPacket.Username == "root" && isPasswordValid {
		// create ok packet
		//fmt.Println("ok packet")
		_, _ = conn.Write(mysql_protocol.NewOKPacket())
	} else {
		// create err packet
		//fmt.Println("err packet")

		var usingPassword = "NO"
		if len(loginPacket.Password) != 0 {
			usingPassword = "YES"
		}

		_, _ = conn.Write(mysql_protocol.NewErrorPacket(fmt.Sprintf("Access denied for user '%s'@'%s' (using password: %s)", loginPacket.Username, conn.RemoteAddr(), usingPassword)))
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	parser, err := sqlparser.New(sqlparser.Options{})
	if err != nil {
		conn.Write(mysql_protocol.NewErrorPacket(err.Error()))
		return
	}

commandLoop:
	for {
		var packetBuffer = make([]byte, s.readBufferSize)
		var _, err = rw.Read(packetBuffer)
		if err != nil {
			// eof
			break
		}

		var pp, parseErr = mysql_protocol.ParseCommandPacket(packetBuffer)
		if parseErr != nil {
			logrus.Fatal(parseErr)
		}
		fmt.Printf("%+#v\n", pp)

		connSession.Query = ""

		switch v := pp.(type) {
		case mysql_protocol.CommandQuery:
			fmt.Println(v.Query)
			var stmt, err = parser.Parse(v.Query)
			if err != nil {
				conn.Write(mysql_protocol.NewErrorPacket(err.Error()))
				continue
			}

			connSession.Query = v.Query

			var stmtResult, stmtResultErr = engine.ExecuteStatement(s.engineConfig, connSession, stmt)
			if stmtResultErr != nil {
				conn.Write(mysql_protocol.NewErrorPacket(err.Error()))
				continue
			}
			fmt.Printf("%+#v\n", stmtResult)

			_, err = conn.Write(mysql_protocol.NewOKPacket())
			if err != nil {
				logrus.Fatal(err)
			}
		case mysql_protocol.CommandQuit:
			break commandLoop
		}

	}
}
