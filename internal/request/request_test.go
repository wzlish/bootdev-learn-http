package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid request-line parts, expecting 3 got 2")

	// Test: Invalid method
	_, err = RequestFromReader(strings.NewReader("INVALID / HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid/unsupported method: INVALID")

	// Test: Empty request target
	_, err = RequestFromReader(strings.NewReader("GET  HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Empty request target path")

	// Test: Unsupported protocol version
	_, err = RequestFromReader(strings.NewReader("GET / HTTP/1.0\r\nHost: localhost:42069\r\n\r\n"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported protocol version, expecting HTTP/1.1 got HTTP/1.0")

	// Test: Request with extra data on the first line (should fail as it won't split into 3 parts)
	_, err = RequestFromReader(strings.NewReader("GET / HTTP/1.1 some extra data\r\nHost: localhost:42069\r\n\r\n"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid request-line parts, expecting 3 got 6")

	// Test: Empty reader input
	_, err = RequestFromReader(strings.NewReader(""))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty reader input")

	// Test: Request with only the request line and no headers
	r, err = RequestFromReader(strings.NewReader("POST /data HTTP/1.1\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/data", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}
