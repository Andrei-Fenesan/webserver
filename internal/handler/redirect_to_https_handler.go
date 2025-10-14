package handler

import (
	"fmt"
	"webserver/internal/model/httpentity"
)

type RedirectToHttpsRequestHandler struct {
	redirectTo string
}

func NewRedirectToHttpsRequestHandler(redirectTo string) *RedirectToHttpsRequestHandler {
	return &RedirectToHttpsRequestHandler{"https://" + redirectTo}
}

func (redirectHandler RedirectToHttpsRequestHandler) ServeRequest(req *httpentity.Request) *httpentity.Response {
	return &httpentity.Response{ResponseCode: 302, Version: req.HttpVersion, Body: nil, Headers: map[string]string{
		"Location":       fmt.Sprintf("%s%s", redirectHandler.redirectTo, req.Path),
		"Content-Type":   "text/html",
		"Content-Length": "0",
	}}
}
