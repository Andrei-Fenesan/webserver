package handler

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	return &httpentity.Response{
		ResponseCode: 200,
		Body:         content,
		Version:      req.HttpVersion,
		Headers: map[string]string{
			"content-type":   getContentType(filepath.Ext(path)),
			"content-length": fmt.Sprintf("%d", len(content)),
		}}

}

func getContentType(fileExtension string) string {
	switch strings.ToLower(fileExtension) {
	case ".html":
		return "text/html"
	case ".txt":
		return "text/plain"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream" // Default for unknown types
	}
}
