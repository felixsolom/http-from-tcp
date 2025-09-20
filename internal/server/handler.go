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
