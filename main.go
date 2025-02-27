package main

import (
	"webserver/internal/handler"
	"webserver/internal/manager"
)

func main() {
	requestHandler := handler.NewHttpRequestHandler("/Users/afenesan/Desktop/personalProj/challenges/webserver/resources")
	connManager := manager.NewConcurrentConnectionManger(requestHandler, 8081)
	connManager.Start()
}
