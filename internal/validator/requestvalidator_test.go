package validator

import (
	"testing"
	"webserver/internal/model/httpentity"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnTrueWhenTheRequestIsValid(t *testing.T) {
	validPaths := []string{"/", "/asd", "/asd/../../", "/cgi-bin/sum.cgi?number=1&number=2&number=3"}
	assert := assert.New(t)

	for _, path := range validPaths {
		req := &httpentity.Request{Path: path, HttpMethod: httpentity.GET}

		assert.True(IsRequestValid(req))
	}
}

func TestShouldReturnFalseWhenTheRequestIsInValid(t *testing.T) {
	invalidPaths := []string{"/asd,", "/`", "/asd/<", "~"}
	assert := assert.New(t)

	for _, path := range invalidPaths {
		req := &httpentity.Request{Path: path, HttpMethod: httpentity.GET}

		assert.False(IsRequestValid(req))
	}
}
