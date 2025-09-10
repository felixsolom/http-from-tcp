package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(req []byte) (RequestLine, error) {
	reqParts := bytes.Split(req, []byte("\r\n"))
	reqLine := reqParts[0]
	reqLineParts := strings.Split(string(reqLine), " ")

	if len(reqLineParts) < 3 {
		return RequestLine{}, fmt.Errorf("Request-line missing parts")
	}

	method := reqLineParts[0]
	for _, char := range method {
		if !unicode.IsUpper(char) {
			return RequestLine{}, fmt.Errorf("Request method is malformed")
		}
	}

	httpVersion := reqLineParts[2]
	httpVersionParts := strings.Split(httpVersion, "/")
	if httpVersionParts[1] != "1.1" {
		return RequestLine{}, fmt.Errorf("Unsupported HTTP version")
	}

	return RequestLine{
		HttpVersion:   httpVersionParts[1],
		RequestTarget: reqLineParts[1],
		Method:        method,
	}, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("coudln't read from HTTP request: %w", err)
	}

	reqLine, err := parseRequestLine(req)
	if err != nil {
		return nil, fmt.Errorf("coudln't parse Request-line: %w", err)
	}

	return &Request{
		RequestLine: reqLine,
	}, nil
}
