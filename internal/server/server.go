package server

import (
	"fmt"
	"net"
)

type Server struct {
	State    string
	Listener net.Listener
}

func Serve(port int) (*Server, error) {
	portStr := fmt.Sprint(port)
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, fmt.Errorf("Coudln't start server on port: %d", port)
	}
	server := Server{
		Listener: l,
		State:    "initialized",
	}
	return &server, nil
}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	s.State = "closed"
	return nil
}
