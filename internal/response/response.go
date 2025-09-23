package response

import (
	"fmt"
	"io"

	"github.com/felixsolom/http-from-tcp/internal/headers"
)

func NewWriter() *Writer {
	return &Writer{
		Headers:     headers.NewHeaders(),
		writerState: writerInitialized,
	}
}

type Writer struct {
	StatusLine  StatusLine
	Headers     headers.Headers
	body        []byte
	writerState writerState
}

type StatusLine struct {
	StatusCode   StatusCode
	ReasonPhrase string
	HttpVersion  string
}

type writerState int

const (
	writerInitialized writerState = iota
	writerWritingHeaders
	writerWritingBody
	writerDone
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerInitialized {
		return fmt.Errorf("Can't write Status-Line in current state")
	}
	switch statusCode {
	case 200:
		w.StatusLine.HttpVersion = "1.1"
		w.StatusLine.StatusCode = 200
		w.StatusLine.ReasonPhrase = "OK"
	case 400:
		w.StatusLine.HttpVersion = "1.1"
		w.StatusLine.StatusCode = 400
		w.StatusLine.ReasonPhrase = "Bad Request"
	case 500:
		w.StatusLine.HttpVersion = "1.1"
		w.StatusLine.StatusCode = 500
		w.StatusLine.ReasonPhrase = "Internal Server Error"
	default:
		w.StatusLine.HttpVersion = "1.1"
		w.StatusLine.StatusCode = 400
		w.StatusLine.ReasonPhrase = "Bad Request"
	}
	w.writerState = writerWritingHeaders
	return nil
}

func GetDefaultHeaders(contentLen int, contentType string) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprint(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = contentType
	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerWritingHeaders {
		return fmt.Errorf("Can't write Headers in current state")
	}
	for key, value := range headers {
		w.Headers[key] = value
	}

	w.writerState = writerWritingBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerWritingBody {
		return 0, fmt.Errorf("Can't write body in current state")
	}
	w.body = append(w.body, p...)
	return len(p), nil
}

func (w *Writer) Flush(writer io.Writer) error {
	if w.writerState == writerDone {
		return fmt.Errorf("Can't flush a response that is already written")
	}
	statusLine := fmt.Sprintf("HTTP/%s %d %s\r\n", w.StatusLine.HttpVersion, w.StatusLine.StatusCode, w.StatusLine.ReasonPhrase)
	if _, err := writer.Write([]byte(statusLine)); err != nil {
		return fmt.Errorf("Couldn'flush status line: %v", err)
	}

	for key, value := range w.Headers {
		header := fmt.Sprintf("%s: %s\r\n", key, value)
		if _, err := writer.Write([]byte(header)); err != nil {
			return fmt.Errorf("Couldn't flush one of the headers: %v", err)
		}
	}

	if _, err := writer.Write([]byte("\r\n")); err != nil {
		return fmt.Errorf("Couldn't flush end header: %v", err)
	}

	if _, err := writer.Write(w.body); err != nil {
		return fmt.Errorf("Couldn't flush body: %v", err)
	}
	w.writerState = writerDone
	return nil
}
