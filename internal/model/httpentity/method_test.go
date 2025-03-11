package httpentity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodFromStringShouldMappedMethods(t *testing.T) {
	assert := assert.New(t)
	methods := map[string]HttpMethod{"get": GET, "post": POST, "put": PUT, "patch": PATCH, "delete": DELETE, "head": HEAD, "options": OPTIONS}
	for key, val := range methods {
		t.Run("Mapped "+key, func(t *testing.T) {
			method, err := methodFromString(key)
			assert.Equal(val, method)
			assert.Nil(err)
		})
	}
}

func TestMethodFromStringShouldReturnErrorWhenUnknownMethod(t *testing.T) {
	assert := assert.New(t)

	_, err := methodFromString("dsa")
	assert.NotNil(err)
}
