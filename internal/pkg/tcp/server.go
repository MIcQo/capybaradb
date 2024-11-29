package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

func NewServer(port uint) *Server {
	return &Server{port: port}
}

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
		conn, err := s.server.Accept()
		if err != nil || conn == nil {
			err = errors.New("could not accept connection")
			break
		}

		go s.handleConnection(conn)
	}
	return
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		req, err := rw.ReadString('\n')
		if err != nil {
			rw.WriteString("failed to read input")
			rw.Flush()
			return
		}

		rw.WriteString(fmt.Sprintf("Request received: %s", req))
		rw.Flush()
	}
}
