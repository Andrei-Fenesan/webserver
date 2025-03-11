package httpentity

import (
	"errors"
	"strings"
)

type HttpMethod uint8

const (
	GET HttpMethod = iota
	POST
	PUT
	PATCH
	DELETE
	HEAD
	OPTIONS
)

func (hm HttpMethod) String() string {
	switch hm {
	case GET:
		return "GET"
	}
	return ""
}

func methodFromString(method string) (HttpMethod, error) {
	switch strings.ToUpper(method) {
	case "GET":
		return GET, nil
	case "POST":
		return POST, nil
	case "PUT":
		return PUT, nil
	case "PATCH":
		return PATCH, nil
	case "DELETE":
		return DELETE, nil
	case "HEAD":
		return HEAD, nil
	case "OPTIONS":
		return OPTIONS, nil
	}
	return 100, errors.New("invalid method")
}
