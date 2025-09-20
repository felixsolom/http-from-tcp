package response

import (
	"fmt"
	"io"

	"github.com/felixsolom/http-from-tcp/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case 200:
		if _, err := w.Write([]byte(
			"HTTP/1.1 200 OK\r\n",
		)); err != nil {
			return err
		}
	case 400:
		if _, err := w.Write([]byte(
			"HTTP/1.1 400 Bad Request\r\n",
		)); err != nil {
			return err
		}
	case 500:
		if _, err := w.Write([]byte(
			"HTTP/1.1 500 Internal Server Error\r\n",
		)); err != nil {
			return err
		}
	default:
		if _, err := w.Write([]byte(
			"\r\n",
		)); err != nil {
			return err
		}
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprint(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		if _, err := w.Write([]byte(
			fmt.Sprintf("%s: %s\r\n", key, value),
		)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("\r\n")); err != nil {
		return err
	}
	return nil
}
