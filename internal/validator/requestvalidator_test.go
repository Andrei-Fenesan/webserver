package validator

import (
	"testing"
	"webserver/internal/model/httpentity"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnTrueWhenTheRequestIsValid(t *testing.T) {
	validPaths := []string{"/", "/asd", "/asd/../../"}
	assert := assert.New(t)

	for _, path := range validPaths {
		req := &httpentity.Request{Path: path, HttpMethod: httpentity.GET}

		assert.True(IsRequestValid(req))
	}
}

func TestShouldReturnFalseWhenTheRequestIsInValid(t *testing.T) {
	validPaths := []string{"/-", "/+", "/asd/<", "~"}
	assert := assert.New(t)

	for _, path := range validPaths {
		req := &httpentity.Request{Path: path, HttpMethod: httpentity.GET}

		assert.False(IsRequestValid(req))
	}
}
