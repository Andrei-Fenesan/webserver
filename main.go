package main

import (
	"flag"
	"log"
	"webserver/internal/handler"
	"webserver/internal/manager"
	"webserver/internal/ssl"
)

func main() {
	var rootDirectory string
	var privateKeyPath string
	var certificatePath string
	var port int
	flag.StringVar(&rootDirectory, "src", "./resources", "The root directory from which the resource will be served")
	flag.IntVar(&port, "port", 80, "The server port")
	flag.StringVar(&certificatePath, "certPath", "", "The path to the certificate file")
	flag.StringVar(&privateKeyPath, "privateKeyPath", "", "The path to the certificate private file")
	flag.Parse()

	actualPort := uint32(port)
	requestHandler := handler.NewHttpRequestHandler(rootDirectory)
	sslHandler := ssl.NewConnectionHandler("./certs/key.pem", "./certs/cert.pem")
	sslHandler.Initialize()
	connManager := manager.NewConcurrentConnectionManger(requestHandler, sslHandler, actualPort)
	log.Printf("The server will start on port: %d, using the root direcotry: %s\n", actualPort, rootDirectory)
	connManager.Start()
}
