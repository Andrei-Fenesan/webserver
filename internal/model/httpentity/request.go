package httpentity

import (
	"errors"
	"fmt"
	"strings"
)

type Request struct {
	HttpVersion string
	HttpMethod  HttpMethod
	Path        string
	Headers     map[string]string
}

func ParseRequest(data []byte) (*Request, error) {
	stringData := string(data)
	lines := strings.Split(stringData, "\r\n")
	if len(lines) < 2 {
		return nil, errors.New("invalid request. no line feed")
	}
	lines = lines[:len(lines)-2] //The last 2 lines are always empty, because the request is finished by \r\n\r\n
	httpMethod, path, version, err := parseStartLine(lines[0])
	if err != nil {
		return nil, err
	}
	headers, err := parseHeaders(lines[1:])
	if err != nil {
		return nil, err
	}

	return &Request{
		version,
		httpMethod,
		path,
		headers,
	}, nil
}

func parseStartLine(startLine string) (HttpMethod, string, string, error) {
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return 0, "", "", errors.New("invalid request")
	}
	method := parts[0]
	path := parts[1]
	version := parts[2]
	httpMethod, err := methodFromString(method)
	if err != nil || !hasValidPathLength(path) || len(version) == 0 {
		return 0, "", "", errors.New("invalid request parts")
	}
	return httpMethod, path, version, nil
}

func parseHeaders(requestHeaders []string) (map[string]string, error) {
	headers := make(map[string]string)
	for _, header := range requestHeaders {
		values := strings.Split(header, ": ")
		if len(values) != 2 {
			return nil, errors.New("invalid headers format")
		}
		headers[values[0]] = values[1]
	}
	return headers, nil
}

func (r *Request) String() string {
	return fmt.Sprintf("%s %s %s\n", r.HttpMethod, r.Path, r.HttpVersion)
}

func hasValidPathLength(path string) bool {
	return len(path) > 0
}
