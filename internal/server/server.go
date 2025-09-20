package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	portStr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, fmt.Errorf("Coudln't start server on port: %d", port)
	}
	server := &Server{
		listener: l,
	}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) listen() {
	for {
		// Wait for a connection.
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				break
			}
			log.Println("Accept error:", err)
			continue
		}
		s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	go func(c net.Conn) {
		c.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/plain\r\n" +
				"\r\n" +
				"Hello World!\n",
		))
		c.Close()
	}(conn)
}
