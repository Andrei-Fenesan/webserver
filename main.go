package main

import (
	"flag"
	"log"
	"sync"
	"webserver/internal/connection-preparer"
	"webserver/internal/handler"
	"webserver/internal/manager"
)

func main() {
	var rootDirectory string
	var privateKeyPath string
	var certificatePath string
	var httpRedirectTo string
	var httpPort int
	var httpsPort int
	flag.StringVar(&rootDirectory, "src", "./resources", "The root directory from which the resource will be served")
	flag.IntVar(&httpPort, "httpPort", 80, "The server httpPort")
	flag.IntVar(&httpsPort, "httpsPort", 443, "The server httpPort")
	flag.StringVar(&certificatePath, "certPath", "", "The path to the certificate file")
	flag.StringVar(&privateKeyPath, "privateKeyPath", "", "The path to the certificate private file")
	flag.StringVar(&httpRedirectTo, "httpRedirectTo", "localhost", "Host where the http traffic should be redirected")
	flag.Parse()

	var serverWaitGroup sync.WaitGroup
	actualHttpPort := uint32(httpPort)
	httpConnManager := initialiseHttpConnectionManager(certificatePath, privateKeyPath, httpRedirectTo, rootDirectory, actualHttpPort)
	serverWaitGroup.Add(1)

	actualHttpsPort := uint32(httpsPort)
	httpsConnManager := initialiseHttpsConnectionManager(certificatePath, privateKeyPath, rootDirectory, actualHttpsPort)

	go func() {
		defer serverWaitGroup.Done()
		log.Printf("The HTTP server will start on port: %d, using the root direcotry: %s\n", actualHttpPort, rootDirectory)
		err := httpConnManager.Start()
		if err != nil {
			log.Println("Error starting HTTP server:", err)
		}
	}()

	if httpsConnManager != nil {
		serverWaitGroup.Add(1)
		go func() {
			defer serverWaitGroup.Done()
			log.Printf("The HTTPS server will start on port: %d, using the root direcotry: %s\n", actualHttpsPort, rootDirectory)
			err := httpsConnManager.Start()
			if err != nil {
				log.Println("Error starting HTTPS server:", err)
			}
		}()
	}
	serverWaitGroup.Wait()
}

func getSslConnectionPreparer(certificatePath string, privateKeyPath string) connection_preparer.ConnectionPreparer {
	log.Println("Loading certificate and private key")
	sslHandler := connection_preparer.NewConnectionHandler(privateKeyPath, certificatePath)
	sslHandler.Initialize()
	return sslHandler
}

func initialiseHttpsConnectionManager(certificatePath string, privateKeyPath string, rootDirectory string, httpsPort uint32) manager.ConnectionManager {
	if certificatePath != "" && privateKeyPath != "" {
		requestHandler := handler.NewHttpRequestHandler(rootDirectory)
		sslConnectionPreparer := getSslConnectionPreparer(certificatePath, privateKeyPath)
		return manager.NewConcurrentConnectionManger(requestHandler, sslConnectionPreparer, httpsPort)
	}
	return nil
}

func initialiseHttpConnectionManager(certificatePath string, privateKeyPath string, redirectTo string, rootDirectory string, httpPort uint32) manager.ConnectionManager {
	sslConnectionPreparer := &connection_preparer.PlainConnectionPreparer{}
	if certificatePath != "" && privateKeyPath != "" {
		requestHandler := handler.NewRedirectToHttpsRequestHandler(redirectTo)
		return manager.NewConcurrentConnectionManger(requestHandler, sslConnectionPreparer, httpPort)
	} else {
		requestHandler := handler.NewHttpRequestHandler(rootDirectory)
		return manager.NewConcurrentConnectionManger(requestHandler, sslConnectionPreparer, httpPort)
	}
}
