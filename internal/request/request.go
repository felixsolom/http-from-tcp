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
	ParserState ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// "Enum" init
type ParserState int

const (
	initialized ParserState = iota
	done
) // End of Enum init

const crlf = "\r\n"
const bufferSize = 8

func (r *Request) parse(data []byte) (int, error) {
	if r.ParserState == 0 {
		reqLine, numOfBytesParsed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if numOfBytesParsed == 0 {
			return 0, nil
		}

		if numOfBytesParsed > 0 {
			r.RequestLine = *reqLine
			r.ParserState = 1
			return numOfBytesParsed, nil
		}
	}
	if r.ParserState == 1 {
		return 0, fmt.Errorf("trying to read data in done state")
	}
	if r.ParserState > 1 {
		return 0, fmt.Errorf("unknown state")
	}
}

func parseRequestLine(req []byte) (*RequestLine, int, error) {
	idx := bytes.Index(req, []byte(crlf))
	if idx == -1 {
		log.Println("No CRLF found. More data needed before request can be parsed")
		return nil, 0, nil
	}
	reqLine := req[:idx]
	parsedReqLine, err := parseRequestLineString(string(reqLine))
	if err != nil {
		return nil, 0, err
	}
	return parsedReqLine, len(req), nil
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
	buff := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	r := Request{
		ParserState: 0,
	}

	for r.ParserState != 1 {
		n, err := reader.Read(buff)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("couldn't read request to buffer: %w", err)
		}
		if err == io.EOF {
			r.ParserState = 1
			break
		}
		readToIndex += n

		if readToIndex > bufferSize {
			newBuffSize := bufferSize * 2
			newBuff := make([]byte, newBuffSize, newBuffSize)
			copy(newBuff, buff)
			buff = newBuff
		}

		numOfBytesParsed, err := r.parse(buff)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse from buffer: %w", err)
		}
		newBuff := make([]byte, numOfBytesParsed)
		copy(newBuff, buff)
		buff = newBuff
		readToIndex -= numOfBytesParsed
	}

	return &r, nil
}
