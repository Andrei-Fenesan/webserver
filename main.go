package main

import (
	"flag"
	"log"
	"webserver/internal/connection-preparer"
	"webserver/internal/handler"
	"webserver/internal/manager"
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
	connectionPreparer := getConnectionPreparer(certificatePath, privateKeyPath)
	connManager := manager.NewConcurrentConnectionManger(requestHandler, connectionPreparer, actualPort)
	log.Printf("The server will start on port: %d, using the root direcotry: %s\n", actualPort, rootDirectory)
	connManager.Start()
}

func getConnectionPreparer(certificatePath string, privateKeyPath string) connection_preparer.ConnectionPreparer {
	if certificatePath != "" && privateKeyPath != "" {
		log.Println("Loading certificate and private key")
		sslHandler := connection_preparer.NewConnectionHandler(privateKeyPath, certificatePath)
		sslHandler.Initialize()
		return sslHandler
	}
	return &connection_preparer.PlainConnectionPreparer{}
}
