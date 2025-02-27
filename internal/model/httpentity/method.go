package httpentity

import (
	"errors"
	"strings"
)

type HttpMethod uint8

const (
	GET HttpMethod = iota
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
	}
	return 100, errors.New("invalid method")
}
