package response

import "io"

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
	}
	return nil
}
