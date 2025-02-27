package handler

import (
	"errors"
	"os"
	"webserver/internal/model/httpentity"
)

type RequestHandler interface {
	ServeRequest(req *httpentity.Request) (*httpentity.Response, error)
}

type HttpRequestHandler struct {
	searchDirectory string
}

func NewHttpRequestHandler(searchDirectory string) *HttpRequestHandler {
	return &HttpRequestHandler{searchDirectory}
}

func (rh *HttpRequestHandler) ServeRequest(req *httpentity.Request) (*httpentity.Response, error) {
	switch req.HttpMethod {
	case httpentity.GET:
		return handleGetRequest(req, rh.searchDirectory), nil
	}
	return nil, errors.New("could not handle the request")
}

func handleGetRequest(req *httpentity.Request, searchDirectpry string) *httpentity.Response {
	path := req.Path
	if path == "/" {
		path = "/index.html"
	}
	content, err := os.ReadFile(searchDirectpry + path)
	if err != nil {
		return &httpentity.Response{ResponseCode: 404, Body: nil, Version: req.HttpVersion}
	}

	return &httpentity.Response{ResponseCode: 200, Body: content, Version: req.HttpVersion}
}
