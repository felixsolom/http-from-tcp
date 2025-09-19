package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Listener net.Listener
	State    bool
}

func Serve(port int) (*Server, error) {
	portStr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, fmt.Errorf("Coudln't start server on port: %d", port)
	}
	server := Server{
		Listener: l,
		State:    true,
	}
	go server.listen()
	return &server, nil
}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	s.State = false
	return nil
}

func (s *Server) listen() {
	for s.State {
		// Wait for a connection.
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.State == false {
				break
			}
			continue
		}
		log.Println("Accept error:", err)
		s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	go func(c net.Conn) {
		c.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/plain\r\n" +
				"Content-Length: 13\r\n" +
				"\r\n" +
				"Hello World!\n",
		))
		// Shut down the connection.
		c.Close()
	}(conn)
}
