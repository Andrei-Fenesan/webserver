package manager

import (
	"io"
	"log"
	"net"
	"path"
	"strconv"
	connection_preparer "webserver/internal/connection-preparer"
	"webserver/internal/handler"
	"webserver/internal/model/httpentity"
	"webserver/internal/validator"
)

const HttpVersion = "HTTP/1.1"
const BuffSize = 1024
const DefaultServerPort = uint32(8080)

type ConnectionManager interface {
	Start() error
	Close()
	handleConnection(conn net.Conn)
}

type ConcurrentConnectionManger struct {
	port               uint32
	requestHandler     handler.RequestHandler
	connectionPreparer connection_preparer.ConnectionPreparer
	listener           net.Listener
}

func NewConcurrentConnectionManger(rq handler.RequestHandler, sh connection_preparer.ConnectionPreparer, port ...uint32) *ConcurrentConnectionManger {
	actualPort := DefaultServerPort
	if len(port) > 0 {
		actualPort = port[0]
	}
	return &ConcurrentConnectionManger{port: actualPort, requestHandler: rq, connectionPreparer: sh}
}

func (cm *ConcurrentConnectionManger) Start() error {
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.FormatUint(uint64(cm.port), 10))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	cm.listener = listener

	for {
		conn, err := listener.Accept()
		go func() {
			if err != nil {
				log.Println("Error in listening " + err.Error())
				return
			}
			defer conn.Close()
			actualConn, err := cm.connectionPreparer.Prepare(conn)
			if err != nil {
				log.Println("Error in preparing" + err.Error())
				return
			}
			cm.handleConnection(actualConn)
		}()
	}
	return nil
}

func (cm *ConcurrentConnectionManger) Close() {
	log.Println("Closing server")
	cm.listener.Close()
}

func (cm *ConcurrentConnectionManger) handleConnection(conn net.Conn) {

	data, err := readAll(conn)
	if err != nil {
		log.Println("Error in reading data " + err.Error())
		return
	}
	log.Println("Received request\n", string(data))

	req, err := httpentity.ParseRequest(data)
	if err != nil {
		log.Println(err)
		conn.Write((&httpentity.Response{ResponseCode: 400, Version: HttpVersion}).Encode())
		return
	}
	if !validator.IsRequestValid(req) {
		log.Printf("Invalid request: %s\n", req.Path)
		conn.Write((&httpentity.Response{ResponseCode: 400, Version: HttpVersion}).Encode())
		return
	}

	canonizeRequest(req)

	response := cm.requestHandler.ServeRequest(req)
	conn.Write(response.Encode())
}

func readAll(conn net.Conn) ([]byte, error) {
	data := make([]byte, 0, 4*BuffSize)
	for {
		buff := make([]byte, BuffSize)
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
	length := len(data)
	if length < 4 {
		return false
	}
	return data[length-4] == '\r' && data[length-3] == '\n' && data[length-2] == '\r' && data[length-1] == '\n'
}

// canonizeRequest will sanitize the request path
func canonizeRequest(req *httpentity.Request) {
	req.Path = path.Clean("/" + req.Path)
}
