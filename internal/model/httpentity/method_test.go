package httpentity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodFromStringShouldReturnGetWhenMethodIsget(t *testing.T) {
	assert := assert.New(t)

	httMethod, err := methodFromString("get")
	assert.Equal(GET, httMethod)
	assert.Nil(err)
}

func TestMethodFromStringShouldReturnErrorWhenUnknownMethod(t *testing.T) {
	assert := assert.New(t)

	_, err := methodFromString("dsa")
	assert.NotNil(err)
}
