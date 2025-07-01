package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(line string) (*Request, error) {

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid request-line parts, expecting 3 got %d", len(parts))
	}

	allowedMethods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}
	if !slices.Contains(allowedMethods, parts[0]) {
		return nil, fmt.Errorf("Invalid/unsupported method: %s", parts[0])
	}

	if parts[1] == "" {
		return nil, fmt.Errorf("Empty request target path")
	}

	if parts[2] != "HTTP/1.1" {
		return nil, fmt.Errorf("Unsupported protocol version, expecting HTTP/1.1 got %s", parts[2])
	}

	return &Request{RequestLine{"1.1", parts[1], parts[0]}}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	readerBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Unable to io.ReadAll: %v", err)
	}
	readerLines := strings.Split(string(readerBytes), "\r\n")

	if len(readerLines) == 0 || (len(readerLines) == 1 && readerLines[0] == "") {
		return nil, fmt.Errorf("empty reader input")
	}

	return parseRequestLine(readerLines[0])
}
