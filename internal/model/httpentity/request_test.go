package httpentity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldParseRequestSuccessfully(t *testing.T) {
	assert := assert.New(t)

	receivedData := "GET /users HTTP/1.1\r\nHost: localhost:8081\r\nUser-Agent: curl/8.7.1\r\nAccept: */*\r\n\r\n"

	request, err := ParseRequest([]byte(receivedData))
	assert.Nil(err)
	assert.NotNil(request)

	assert.Equal(GET, request.HttpMethod)
	assert.Equal("/users", request.Path)
	assert.Equal("HTTP/1.1", request.HttpVersion)
}

func TestShouldParseRequestHeadersSuccessfully(t *testing.T) {
	assert := assert.New(t)

	receivedData := "GET /users HTTP/1.1\r\nHost: localhost:8081\r\nAccept: */*\r\n\r\n"

	request, err := ParseRequest([]byte(receivedData))
	assert.Nil(err)
	assert.NotNil(request)

	assert.Equal(2, len(request.Headers))
	assert.Equal("localhost:8081", request.Headers["Host"])
	assert.Equal("*/*", request.Headers["Accept"])
}

func TestShouldParseRequestHeadersFailsWhenInvalidHeaders(t *testing.T) {
	assert := assert.New(t)

	receivedData := "GET /users HTTP/1.1\r\nHost-localhost:8081\r\n\r\n"

	request, err := ParseRequest([]byte(receivedData))
	assert.Nil(request)
	assert.NotNil(err)
}

func TestShouldParseRequestShouldFaildWhenUnknownHttpMethod(t *testing.T) {
	assert := assert.New(t)

	receivedData := "ASD /users HTTP/1.1\r\nHost: localhost:8081\r\nUser-Agent: curl/8.7.1\r\nAccept: */*\r\n\r\n"

	request, err := ParseRequest([]byte(receivedData))
	assert.NotNil(err)
	assert.Nil(request)
}

func TestShouldParseRequestShouldFaildWhenNoPath(t *testing.T) {
	assert := assert.New(t)

	receivedData := "GET HTTP/1.1\r\nHost: localhost:8081\r\nUser-Agent: curl/8.7.1\r\nAccept: */*\r\n\r\n"

	request, err := ParseRequest([]byte(receivedData))
	assert.NotNil(err)
	assert.Nil(request)
}

func TestShouldParseRequestShouldFaildWhenNoVersion(t *testing.T) {
	assert := assert.New(t)

	receivedData := "GET /asd\r\nHost: localhost:8081\r\nUser-Agent: curl/8.7.1\r\nAccept: */*\r\n\r\n"

	request, err := ParseRequest([]byte(receivedData))
	assert.NotNil(err)
	assert.Nil(request)
}
