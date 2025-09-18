package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid end of headers
	headers = NewHeaders()
	data = []byte("\r\n things more things")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Valid 2 headers with one already existing in map
	headers = map[string]string{"host": "localhost:42069"}
	data = []byte("Content-Type: application/json; charset=utf-8\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "application/json; charset=utf-8", headers["content-type"])
	assert.Equal(t, 47, n)
	assert.False(t, done)

	// Test: Valid 2 headers with the same field-name, first already existing in map
	headers = map[string]string{"set-person": "lane-loves-go"}
	data = []byte("Set-Person: prime-loves-zig\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "lane-loves-go, prime-loves-zig", headers["set-person"])
	assert.Equal(t, 29, n)
	assert.False(t, done)

	// Test: Invalid single header with invalid chars in field-name
	data = []byte("Content-T@pe: application/json; charset=utf-8\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

}
