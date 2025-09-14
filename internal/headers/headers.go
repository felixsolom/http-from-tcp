package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return len(crlf), true, nil
	}
	newHeader := string(data[:idx])
	colonIdx := strings.Index(newHeader, ":")
	if colonIdx == 0 || colonIdx == -1 || newHeader[colonIdx-1] == ' ' {
		return 0, false, fmt.Errorf("Malformed header")
	}
	fieldLine := strings.SplitN(newHeader, ":", 2)
	if len(fieldLine) != 2 {
		return 0, false, fmt.Errorf("Malformed header: %v", fieldLine)
	}

	fieldName := strings.TrimSpace(fieldLine[0])
	fieldValue := strings.TrimSpace(fieldLine[1])
	h[fieldName] = fieldValue
	return idx + len(crlf), false, nil
}
