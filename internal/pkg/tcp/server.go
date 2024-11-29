// Package tcp contains the TCP Server
package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

// NewServer creates a new TCP Server
func NewServer(port uint) *Server {
	return &Server{port: port}
}

// Server represents the TCP Server
type Server struct {
	port   uint
	server net.Listener
}

// Start starts the TCP Server
func (s *Server) Start() (err error) {

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

		go s.handleConnection(conn)
	}
	return
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		req, err := rw.ReadString('\n')
		if err != nil {
			_, _ = rw.WriteString("failed to read input")
			_ = rw.Flush()
			return
		}

		_, _ = rw.WriteString(fmt.Sprintf("Request received: %s", req))
		_ = rw.Flush()
	}
}
