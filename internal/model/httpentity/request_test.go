package httpentity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldParseRequestSuccessfully(t *testing.T) {
	assert := assert.New(t)

	receivedData := `GET /users HTTP/1.1
					Host: localhost:8081
					User-Agent: curl/8.7.1
					Accept: */*`

	request, err := ParseRequest([]byte(receivedData))
	assert.Nil(err)
	assert.NotNil(request)

	assert.Equal(GET, request.HttpMethod)
	assert.Equal("/users", request.Path)
	assert.Equal("HTTP/1.1", request.HttpVersion)
}

func TestShouldParseRequestShouldFaildWhenUnknownHttpMethod(t *testing.T) {
	assert := assert.New(t)

	receivedData := `ASD /users HTTP/1.1
					Host: localhost:8081
					User-Agent: curl/8.7.1
					Accept: */*`

	request, err := ParseRequest([]byte(receivedData))
	assert.NotNil(err)
	assert.Nil(request)
}

func TestShouldParseRequestShouldFaildWhenNoPath(t *testing.T) {
	assert := assert.New(t)

	receivedData := `GET HTTP/1.1
					Host: localhost:8081
					User-Agent: curl/8.7.1
					Accept: */*`

	request, err := ParseRequest([]byte(receivedData))
	assert.NotNil(err)
	assert.Nil(request)
}

func TestShouldParseRequestShouldFaildWhenNoVersion(t *testing.T) {
	assert := assert.New(t)

	receivedData := `GET /asd
					Host: localhost:8081
					User-Agent: curl/8.7.1
					Accept: */*`

	request, err := ParseRequest([]byte(receivedData))
	assert.NotNil(err)
	assert.Nil(request)
}

func TestShouldParseRequestShouldFaildWhenNoFinal(t *testing.T) {
	assert := assert.New(t)

	receivedData := `GET /asd Host: localhost:8081 User-Agent: curl/8.7.1 Accept: */*`

	request, err := ParseRequest([]byte(receivedData))
	assert.NotNil(err)
	assert.Nil(request)
}
