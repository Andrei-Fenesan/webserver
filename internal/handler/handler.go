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
	"webserver/internal/linetermination"
	"webserver/internal/model/httpentity"
)

const PathInfoEnvName = "PATH_INFO"
const QueryStringEnvName = "QUERY_STRING"

type RequestHandler interface {
	ServeRequest(req *httpentity.Request) *httpentity.Response
}

type HttpRequestHandler struct {
	searchDirectory string
}

func NewHttpRequestHandler(searchDirectory string) *HttpRequestHandler {
	return &HttpRequestHandler{searchDirectory}
}

func (rh *HttpRequestHandler) ServeRequest(req *httpentity.Request) *httpentity.Response {
	log.Printf("Serving the request: %s\n", req)
	switch req.HttpMethod {
	case httpentity.GET:
		return handleGetRequest(req, rh.searchDirectory)
	default:
		return &httpentity.Response{ResponseCode: 405, Body: nil, Version: req.HttpVersion}
	}
}

func handleGetRequest(req *httpentity.Request, searchDirectory string) *httpentity.Response {
	path := req.Path
	if path == "/" {
		path = "/index.html"
	}
	searchPath := searchDirectory + path
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
	cgiExecutablePath, queryParams, err := parseSearchPath(searchPath)
	if err != nil {
		log.Printf("Cgi not found in search path: %s\n", searchPath)
		return &httpentity.Response{ResponseCode: 404, Body: nil, Version: req.HttpVersion}
	}

	cmd := exec.Command(cgiExecutablePath)
	cmd.Env = prepareCgiEnvs(queryParams)
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

func parseSearchPath(searchPath string) (cgiExecutablePath string, queryParams string, err error) {
	const cgiExtensionLength = 4
	cgiIndex := strings.Index(searchPath, ".cgi")
	if cgiIndex == -1 {
		//no cgi executable found
		return "", "", errors.New("cgi not found")
	}

	cgiExecutablePath = searchPath[:(cgiIndex + cgiExtensionLength)]
	queryParams = searchPath[(cgiIndex + cgiExtensionLength):]
	err = nil
	return
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
		createEnv(PathInfoEnvName, pathInfoEnvValue),
		createEnv(QueryStringEnvName, queryEnvValue),
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
	returnedHeaders, after, _ := bytes.Cut(body, linetermination.GetLineTermination())
	stringReturnedHeaders := string(returnedHeaders)
	for _, headerLine := range strings.Split(stringReturnedHeaders, string([]byte{'\n'})) {
		headerSplit := strings.Split(headerLine, ": ")
		headers[headerSplit[0]] = headerSplit[1]
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
