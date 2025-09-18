package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
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
		// the empty line
		// headers are done, we are consuming the CRLF
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
	if !isValidFieldName(fieldName) {
		return 0, false, fmt.Errorf("Malformed field-name: %s", fieldName)
	}
	fieldValue := strings.TrimSpace(fieldLine[1])

	//if field-name already exists in a map, we append the new value to existing one
	if existingFieldValue, exists := h[strings.ToLower(fieldName)]; exists {
		h[strings.ToLower(fieldName)] = existingFieldValue + ", " + fieldValue
	} else {
		h[strings.ToLower(fieldName)] = fieldValue
	}
	return idx + len(crlf), false, nil
}

func isValidFieldName(fname string) bool {
	specialChars := "!#$%&'*+-.^_`|~"

	if len(fname) < 1 {
		return false
	}

	for _, char := range fname {
		if !unicode.IsLetter(char) &&
			!unicode.IsDigit(char) &&
			!strings.ContainsRune(specialChars, char) {
			return false
		}
	}
	return true
}
