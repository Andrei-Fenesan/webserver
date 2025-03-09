package internal

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"testing"
	"webserver/internal/handler"
	"webserver/internal/manager"
)

func buildTestServer() manager.ConnectionManager {
	requestHandler := handler.NewHttpRequestHandler("../resources_test")
	return manager.NewConcurrentConnectionManger(requestHandler, 9191)
}

func getFileContent(pagePath string) []byte {
	path := "../resources_test/" + pagePath
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}

func TestWebServer(t *testing.T) {
	tt := []struct {
		testName         string
		payload          []byte
		expectedResponse []byte
	}{
		//200 OK requests
		{
			"Should return 200 OK and /index.html page from root server directory",
			[]byte("GET / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: testagent\r\nAccept: */*\r\n\r\n"),
			append([]byte("HTTP/1.1 200 OK\r\ncontent-type: text/html\r\ncontent-length: 183\r\n\r\n"), getFileContent("index.html")...),
		},
		{
			"Should get 200 OK and /simple.html page from other directory inside the server root directory",
			[]byte("GET /deep_pages/simple.html HTTP/1.1\r\nHost: localhost\r\nUser-Agent: testagent\r\nAccept: */*\r\n\r\n"),
			append([]byte("HTTP/1.1 200 OK\r\ncontent-type: text/html\r\ncontent-length: 419\r\n\r\n"), getFileContent("deep_pages/simple.html")...),
		},
		//400 BAD REQUEST requests
		{
			"Should get 400 BAD REQUEST when method is not GET",
			[]byte("POST /unknown.html HTTP/1.1\r\nHost: localhost\r\nUser-Agent: testagent\r\nAccept: */*\r\n\r\n"),
			[]byte("HTTP/1.1 400 BAD REQUEST\r\n\r\n"),
		},
		{
			"Should get 400 BAD REQUEST when the first line is malformated. Method is missing",
			[]byte("/unknown.html HTTP/1.1\r\nHost: localhost\r\nUser-Agent: testagent\r\nAccept: */*\r\n\r\n"),
			[]byte("HTTP/1.1 400 BAD REQUEST\r\n\r\n"),
		},
		{
			"Should get 400 BAD REQUEST when the first line is malformated. Path is missing",
			[]byte("GET HTTP/1.1\r\nHost: localhost\r\nUser-Agent: testagent\r\nAccept: */*\r\n\r\n"),
			[]byte("HTTP/1.1 400 BAD REQUEST\r\n\r\n"),
		},
		{
			"Should get 400 BAD REQUEST when the first line is malformated. More parts than allowed",
			[]byte("GET / HTTP/1.1 asd\r\nHost: localhost\r\nUser-Agent: testagent\r\nAccept: */*\r\n\r\n"),
			[]byte("HTTP/1.1 400 BAD REQUEST\r\n\r\n"),
		},
		{
			"Should get 400 BAD REQUEST when the headers are malformated",
			[]byte("GET HTTP/1.1\r\nHost-localhost\r\nUser-Agent: testagent\r\nAccept: */*\r\n\r\n"),
			[]byte("HTTP/1.1 400 BAD REQUEST\r\n\r\n"),
		},
		//404 NOT FOUND requests
		{
			"Should get 404 NOT FOUND when searching for a file that does not exist",
			[]byte("GET /unknown.html HTTP/1.1\r\nHost: localhost\r\nUser-Agent: testagent\r\nAccept: */*\r\n\r\n"),
			[]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"),
		},
		{
			"Should get 404 NOT FOUND when searching for a file that does not exist. Should protect from path traversal atack vector",
			[]byte("GET ../main.go HTTP/1.1\r\nHost: localhost\r\nUser-Agent: testagent\r\nAccept: */*\r\n\r\n"),
			[]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"),
		},
	}

	testServer := buildTestServer()
	go testServer.Start()
	defer testServer.Close()

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			conn, err := net.Dial("tcp", "localhost:9191")
			if err != nil {
				t.Error("could not connect to TCP server: ", err)
			}
			defer conn.Close()

			if _, err := conn.Write(tc.payload); err != nil {
				t.Error("could not write payload to TCP server:", err)
			}

			response := make([]byte, 4056)
			if read, err := conn.Read(response); err == nil {
				if !bytes.Equal(response[:read], tc.expectedResponse) {
					fmt.Printf("response[:read]: %s\n", string(response[:read]))
					t.Error("response did not match expected output")
				}
			} else {
				t.Error("could not read from connection")
			}
		})
	}
}
