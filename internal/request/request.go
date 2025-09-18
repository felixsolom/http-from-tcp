package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"

	"github.com/felixsolom/http-from-tcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
	Headers     headers.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// "Enum" init
type ParserState int

const (
	stateInitialized ParserState = iota
	stateParsingHeaders
	stateDone
) // End of Enum init

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize)
	readToIndex := 0

	r := Request{
		ParserState: stateInitialized,
		Headers:     headers.NewHeaders(),
	}

	for r.ParserState != stateDone {
		if readToIndex == len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		numOfBytesRead, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if r.ParserState != stateDone {
					return nil, fmt.Errorf("Incomplete request, in %d, read n bytes on EOF: %d", r.ParserState, numOfBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += numOfBytesRead

		if err == io.EOF && readToIndex == 0 {
			r.ParserState = stateDone
			break
		}

		numOfBytesParsed, parseErr := r.parse(buff[:readToIndex])
		if parseErr != nil {
			return nil, fmt.Errorf("couldn't parse from buffer: %w", parseErr)
		}

		// Shifting the yet unparsed data to the beginning of the buffer.
		if numOfBytesParsed > 0 {
			copy(buff, buff[numOfBytesParsed:readToIndex])
			readToIndex -= numOfBytesParsed
		}

		if err == io.EOF {
			break
		}
	}
	return &r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.ParserState != stateDone {
		numOfBytesParsed, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, fmt.Errorf("couldn't parse headers: %w", err)
		}

		totalBytesParsed += numOfBytesParsed
		if numOfBytesParsed == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.ParserState {
	case stateInitialized:
		reqLine, numOfBytesParsed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if numOfBytesParsed == 0 {
			//not enough data, waiting for more
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.ParserState = stateParsingHeaders
		return numOfBytesParsed, nil

	case stateParsingHeaders:
		numOfBytesParsed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, fmt.Errorf("Couldn't parse headers: %w", err)
		}

		if done {
			r.ParserState = stateDone
		}
		return numOfBytesParsed, nil
	case stateDone:
		return 0, fmt.Errorf("Trying to read data in Done state")
	default:
		return 0, fmt.Errorf("Unknown state")
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
	return parsedReqLine, idx + len(crlf), nil
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
