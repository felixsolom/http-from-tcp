package response

import (
	"fmt"
	"io"

	"github.com/felixsolom/http-from-tcp/internal/headers"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

type Writer struct {
	writerState writerState
	writer      io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer:      w,
		writerState: writerStateStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("cannot write status line in state: %d", w.writerState)
	}
	defer func() { w.writerState = writerStateHeaders }()
	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return fmt.Errorf("cannot write headers in state: %d", w.writerState)
	}

	defer func() { w.writerState = writerStateBody }()
	for key, value := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return nil
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in current state: %d", w.writerState)
	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in current state: %d", w.writerState)
	}
	if len(p) > 0 {
		n, err := w.writer.Write([]byte(
			fmt.Sprintf("%x\r\n%s\r\n", len(p), p)))
		if err != nil {
			return 0, err
		}
		return n, nil
	}
	return 0, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in current state: %d", w.writerState)
	}
	n, err := w.writer.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return 0, err
	}
	return n, nil
}
