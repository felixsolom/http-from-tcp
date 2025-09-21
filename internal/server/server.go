package server

import (
	"bytes"
	"fmt"
	"io"
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

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(messageBytes)))
	w.Write(messageBytes)
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
		//parsing request from the connection
		req, err := request.RequestFromReader(c)
		if err != nil {
			hErr := &HandlerError{
				StatusCode: response.BadRequest,
				Message:    err.Error(),
			}
			hErr.Write(c)
			return
		}
		//writing response
		buf := bytes.NewBuffer([]byte{})
		hErr := s.handler(buf, req)
		if hErr != nil {
			hErr.Write(c)
			return
		}
		b := buf.Bytes()
		response.WriteStatusLine(c, response.OK)
		response.WriteHeaders(c, response.GetDefaultHeaders(len(b)))
		c.Write(b)
		c.Close()
	}(conn)
}
