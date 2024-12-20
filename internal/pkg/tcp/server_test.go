package tcp

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

var srv *Server

const port = 2121

func init() {
	// Start the new server.
	srv = NewServer(WithPort(port))

	// Run the server in Goroutine to stop tests from blocking
	// test execution.
	go func() {
		_ = srv.Start()
	}()
}

func TestServer_Run(t *testing.T) {
	// Simply check that the server is up and can
	// accept connections.
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
}

func TestServer_Request(t *testing.T) {
	tt := []struct {
		test    string
		payload []byte
		want    []byte
	}{
		{
			"Sending a simple request returns result",
			[]byte("hello world\n"),
			[]byte("Request received: hello world"),
		},
		{
			"Sending another simple request works",
			[]byte("goodbye world\n"),
			[]byte("Request received: goodbye world"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.test, func(t *testing.T) {
			conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				t.Error("could not connect to TCP server: ", err)
			}
			defer func(conn net.Conn) {
				_ = conn.Close()
			}(conn)

			if _, err := conn.Write(tc.payload); err != nil {
				t.Error("could not write payload to TCP server:", err)
			}

			out := make([]byte, 1024)
			if _, err := conn.Read(out); err == nil {
				if bytes.Equal(out, tc.want) {
					t.Error("response did match expected output")
				}
			} else {
				t.Error("could not read from connection")
			}
		})
	}
}
