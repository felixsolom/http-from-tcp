package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/felixsolom/http-from-tcp/internal/request"
	"github.com/felixsolom/http-from-tcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	portStr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, fmt.Errorf("Coudln't start server on port: %d", port)
	}
	server := &Server{
		listener: l,
		handler:  handler,
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
		defer c.Close()
		w := response.NewWriter(c)
		req, err := request.RequestFromReader(c)
		if err != nil {
			w.WriteStatusLine(response.BadRequest)
			body := []byte(fmt.Sprintf("error parsing request: %v", err))
			w.WriteHeaders(response.GetDefaultHeaders(len(body)))
			w.WriteBody(body)
			return
		}
		s.handler(w, req)
	}(conn)
}
