package handler

import (
	"testing"
	"webserver/internal/model/httpentity"

	"github.com/stretchr/testify/assert"
)

func TestShouldRedirectToHttpsLocation(t *testing.T) {
	reqHandler := RedirectToHttpsRequestHandler{redirectTo: "https://host.com"}
	response := reqHandler.ServeRequest(buildRequest())
	assert.Equal(t, "https://host.com/pages/index.html", response.Headers["Location"])
}

func buildRequest() *httpentity.Request {
	return &httpentity.Request{
		HttpVersion: "1.1",
		HttpMethod:  httpentity.GET,
		Path:        "/pages/index.html",
		Headers:     map[string]string{},
	}
}
