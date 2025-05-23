package httpentity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnTheCorrectHttpStatusCodeMessage(t *testing.T) {
	assert := assert.New(t)
	codesAndMessages := map[uint16]string{200: "OK", 400: "BAD REQUEST", 404: "NOT FOUND", 405: "METHOD NOT ALLOWED", 500: "INTERNAL SERVER ERROR"}

	for code, message := range codesAndMessages {
		assert.Equal(message, getResponseMessage(code))
	}
}
