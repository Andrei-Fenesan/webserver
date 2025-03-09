package handler

import (
	"errors"
	"log"
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
	log.Printf("Serving the request: %s\n", req)
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
	searchPath := searchDirectpry + path
	log.Printf("Seaching for file: %s\n", searchPath)
	content, err := os.ReadFile(searchPath)
	if err != nil {
		return &httpentity.Response{ResponseCode: 404, Body: nil, Version: req.HttpVersion}
	}

	return &httpentity.Response{ResponseCode: 200, Body: content, Version: req.HttpVersion}
}
