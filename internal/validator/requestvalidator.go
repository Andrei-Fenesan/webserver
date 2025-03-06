package validator

import (
	"regexp"
	"webserver/internal/model/httpentity"
)

var validPathRegex = regexp.MustCompile("^[a-zA-Z0-9$_@!%^&*()./]*$")

func IsRequestValid(req *httpentity.Request) bool {
	return isPathValid(req.Path)
}

func isPathValid(path string) bool {
	return containsValidCharacters(path)
}

func containsValidCharacters(path string) bool {
	return validPathRegex.MatchString(path)
}
