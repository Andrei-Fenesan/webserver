package handler

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"webserver/internal/model/httpentity"
)

const PATH_INFO_ENV_NAME = "PATH_INFO"
const QUERY_STRING_ENV_NAME = "QUERY_STRING"

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
	if isCGI(searchPath) {
		return handleCGIRequest(req, searchPath)
	} else {
		return handleSimpleRequest(req, searchPath)
	}
}

func isCGI(path string) bool {
	return strings.Contains(path, "/cgi-bin")
}

func handleCGIRequest(req *httpentity.Request, searchPath string) *httpentity.Response {
	//search for the CGI exec and add PATH_INFO and QUERY_PARAM variables
	cgiIndex := strings.Index(searchPath, ".cgi")
	if cgiIndex == -1 {
		//no cgi executable found
		return &httpentity.Response{ResponseCode: 404, Body: nil, Version: req.HttpVersion}
	}
	cgiExecutablePath := searchPath[:(cgiIndex + 4)]

	cmd := exec.Command(cgiExecutablePath)
	cmd.Env = prepareCgiEnvs(searchPath[(cgiIndex + 4):])
	data, err := cmd.Output()
	if err == nil {
		// extract headers and add the content-length header
		headers, body := extractHeadersAndBody(data)
		headers["content-length"] = fmt.Sprintf("%d", len(body))
		return &httpentity.Response{
			ResponseCode: 200,
			Body:         body,
			Version:      req.HttpVersion,
			Headers:      headers,
		}
	}
	return &httpentity.Response{ResponseCode: 404, Body: nil, Version: req.HttpVersion}
}

func prepareCgiEnvs(path string) []string {
	var pathInfoEnvValue string
	var queryEnvValue string
	queryIndex := strings.LastIndexByte(path, '?')
	if queryIndex == -1 {
		pathInfoEnvValue = path
	} else {
		pathInfoEnvValue = path[:queryIndex]
		queryEnvValue = path[queryIndex:]
	}
	return []string{
		createEnv(PATH_INFO_ENV_NAME, pathInfoEnvValue),
		createEnv(QUERY_STRING_ENV_NAME, queryEnvValue),
	}
}

func createEnv(key string, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

func handleSimpleRequest(req *httpentity.Request, searchPath string) *httpentity.Response {
	content, err := os.ReadFile(searchPath)
	if err != nil {
		return &httpentity.Response{ResponseCode: 404, Body: nil, Version: req.HttpVersion}
	}

	return &httpentity.Response{
		ResponseCode: 200,
		Body:         content,
		Version:      req.HttpVersion,
		Headers: map[string]string{
			"content-type":   getContentType(filepath.Ext(searchPath)),
			"content-length": fmt.Sprintf("%d", len(content)),
		}}
}

func extractHeadersAndBody(body []byte) (map[string]string, []byte) {
	headers := make(map[string]string)
	returedHeaders, after, _ := bytes.Cut(body, []byte{'\n', '\n'})
	stringReturnedHeaders := string(returedHeaders)
	for _, headerLine := range strings.Split(stringReturnedHeaders, string([]byte{'\n'})) {
		headerSplited := strings.Split(headerLine, ": ")
		headers[headerSplited[0]] = headerSplited[1]
	}
	return headers, after
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
