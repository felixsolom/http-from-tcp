package server

import (
	"io"

	"github.com/felixsolom/http-from-tcp/internal/request"
	"github.com/felixsolom/http-from-tcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (h Handler) WriteHandlerError(w io.Writer, message string, status response.StatusCode) *HandlerError {
	return &HandlerError{
		StatusCode: status,
		Message:    message,
	}
}
