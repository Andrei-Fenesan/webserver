package main

import (
	"flag"
	"log"
	"webserver/internal/handler"
	"webserver/internal/manager"
)

func main() {
	var rootDirectory string
	var port int
	flag.StringVar(&rootDirectory, "src", "./resources", "The root directory from which the resource will be served")
	flag.IntVar(&port, "port", 8080, "The server port")
	flag.Parse()

	actualPort := uint32(port)
	requestHandler := handler.NewHttpRequestHandler(rootDirectory)
	connManager := manager.NewConcurrentConnectionManger(requestHandler, actualPort)
	log.Printf("The server will start on port: %d, using the root direcotry: %s\n", actualPort, rootDirectory)
	connManager.Start()
}
