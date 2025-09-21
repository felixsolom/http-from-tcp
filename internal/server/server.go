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
			log.Printf("Couldn't get request-line from reader: %v", err)

			if err := response.WriteStatusLine(c, response.BadRequest); err != nil {
				log.Printf("Couldn't write status line: %v", err)
			}
			if err := response.WriteHeaders(c, response.GetDefaultHeaders(0)); err != nil {
				log.Printf("Couldn't write headers: %v", err)
			}
		}
		//writing response
		if handlerErr := s.handler(c, req); handlerErr != nil {
			log.Printf("Handler error: %s", handlerErr.Message)

			if err := response.WriteStatusLine(c, handlerErr.StatusCode); err != nil {
				log.Printf("Couldn't write status line: %v", err)
			}
			body := handlerErr.Message
			headers := response.GetDefaultHeaders(len(body))
			if err := response.WriteHeaders(c, headers); err != nil {
				log.Printf("Couldn't write headers: %v", err)
			}
			if _, err := c.Write([]byte(body)); err != nil {
				log.Printf("Coudln't write response body: %v", err)
			}
		}
		c.Close()
	}(conn)
}
