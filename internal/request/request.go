package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
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

func parseRequestLine(req []byte) (*RequestLine, int, error) {
	idx := bytes.Index(req, []byte("\r\n"))
	if idx == -1 {
		log.Println("No SLRF found. More data needed before request can be parsed")
		return nil, 0, nil
	}
	reqLine := req[:idx]
	parsedReqLine, err := parseRequestLineString(string(reqLine))
	if err != nil {
		return nil, 0, err
	}
	return parsedReqLine, 0, nil
}

func parseRequestLineString(reqLine string) (*RequestLine, error) {
	reqLineParts := strings.Split(string(reqLine), " ")
	if len(reqLineParts) < 3 {
		return nil, fmt.Errorf("request-line missing parts: %v", reqLineParts)
	}

	method := reqLineParts[0]
	for _, char := range method {
		if !unicode.IsUpper(char) {
			return nil, fmt.Errorf("request method is malformed: %s", method)
		}
	}

	httpVersion := reqLineParts[2]

	httpVersionParts := strings.Split(httpVersion, "/")
	if len(httpVersionParts) < 2 {
		return nil, fmt.Errorf("malformed HTTP version: %s", httpVersionParts)
	}
	if httpVersionParts[0] != "HTTP" {
		return nil, fmt.Errorf("malformed HTTP version: %s", httpVersionParts[0])
	}
	if httpVersionParts[1] != "1.1" {
		return nil, fmt.Errorf("unsupported HTTP version: %s", httpVersionParts[1])
	}

	return &RequestLine{
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

	reqLine, _, err := parseRequestLine(req)
	if err != nil {
		return nil, fmt.Errorf("coudln't parse Request-line: %w", err)
	}

	return &Request{
		RequestLine: *reqLine,
	}, nil
}
