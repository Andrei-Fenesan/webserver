package manager

import (
	"io"
	"log"
	"net"
	"path/filepath"
	"strconv"
	"webserver/internal/handler"
	"webserver/internal/model/httpentity"
	"webserver/internal/validator"
)

const HTTP_VERSION = "HTTP/1.1"
const BUFF_SIZE = 1024
const DEFAULT_SERVER_PORT = uint32(8080)

type ConnectionManager interface {
	Start() error
	Close()
	handleConnection(conn net.Conn)
}

type ConcurrentConnectionManger struct {
	port           uint32
	requestHandler handler.RequestHandler
	listener       net.Listener
}

func NewConcurrentConnectionManger(rq handler.RequestHandler, port ...uint32) *ConcurrentConnectionManger {
	actualPort := DEFAULT_SERVER_PORT
	if len(port) > 0 {
		actualPort = port[0]
	}
	return &ConcurrentConnectionManger{port: actualPort, requestHandler: rq}
}

func (cm *ConcurrentConnectionManger) Start() error {
	listener, err := net.Listen("tcp", "localhost:"+strconv.FormatUint(uint64(cm.port), 10))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	cm.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error in listening" + err.Error())
			break
		}
		go cm.handleConnection(conn)
	}
	return nil
}

func (cm *ConcurrentConnectionManger) Close() {
	log.Println("Closing server")
	cm.listener.Close()
}

func (cm *ConcurrentConnectionManger) handleConnection(conn net.Conn) {
	defer conn.Close()

	data, err := readAll(conn)
	if err != nil {
		log.Println("Error in listening" + err.Error())
		return
	}
	log.Println("Received request\n", string(data))

	req, err := httpentity.ParseRequest(data)
	if err != nil {
		log.Println(err)
		conn.Write((&httpentity.Response{ResponseCode: 400, Version: HTTP_VERSION}).Encode())
		return
	}
	if !validator.IsRequestValid(req) {
		log.Printf("Invalid request: %s\n", req.Path)
		conn.Write((&httpentity.Response{ResponseCode: 400, Version: HTTP_VERSION}).Encode())
		return
	}

	canonizeRequest(req)

	response, err := cm.requestHandler.ServeRequest(req)
	if err != nil {
		log.Println(err)
		return
	}
	conn.Write(response.Encode())
}

func readAll(conn net.Conn) ([]byte, error) {
	data := make([]byte, 0, 4*BUFF_SIZE)
	for {
		buff := make([]byte, BUFF_SIZE)
		read, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return nil, err
		}
		buff = buff[:read]
		data = append(data, buff...)
		if isReadingFinished(data) {
			break
		}
	}
	return data, nil
}

func isReadingFinished(data []byte) bool {
	len := len(data)
	if len < 4 {
		return false
	}
	return data[len-4] == '\r' && data[len-3] == '\n' && data[len-2] == '\r' && data[len-1] == '\n'
}

// canonizeRequest will sanitize the request path
func canonizeRequest(req *httpentity.Request) {
	req.Path = filepath.Clean(req.Path)
}
