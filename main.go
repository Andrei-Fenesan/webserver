package main

import (
	"flag"
	"log"
	"webserver/internal/handler"
	"webserver/internal/manager"
)

func main() {
	var rootDirectory string
	flag.StringVar(&rootDirectory, "src", "./resources", "The root direcotory from which the resource will be served")
	flag.Parse()

	requestHandler := handler.NewHttpRequestHandler(rootDirectory)
	connManager := manager.NewConcurrentConnectionManger(requestHandler, 8081)
	log.Printf("The server will start using the root direcotry: %s\n", rootDirectory)
	connManager.Start()
}
