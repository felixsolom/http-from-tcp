package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/felixsolom/http-from-tcp/internal/request"
	"github.com/felixsolom/http-from-tcp/internal/response"
)

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
		//parsing request from the connection
		req, err := request.RequestFromReader(c)
		if err != nil {
			rw := response.NewWriter()
			rw.WriteStatusLine(response.BadRequest)
			rw.WriteHeaders(response.GetDefaultHeaders(0, "text/html"))
			if err := rw.Flush(c); err != nil {
				log.Printf("Couldn't flush error response: %v", err)
			}
			return
		}

		//writing response
		rw := response.NewWriter()
		if hErr := s.handler(rw, req); hErr != nil {
			log.Printf("handler error: %s", hErr.Message)
			errorRw := response.NewWriter()
			errorRw.WriteStatusLine(hErr.StatusCode)
			body := []byte(hErr.Message)
			errorRw.WriteHeaders(response.GetDefaultHeaders(len(body), "text/html"))
			errorRw.WriteBody(body)
			if err := errorRw.Flush(c); err != nil {
				log.Printf("Couldn't flush error response: %v", err)
			}
			return
		}

		if err := rw.Flush(c); err != nil {
			log.Printf("Couldn't flush response: %v", err)
		}
	}(conn)
}
