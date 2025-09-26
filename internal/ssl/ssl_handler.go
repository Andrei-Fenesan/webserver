package ssl

import (
	"crypto/tls"
	"log"
	"net"
)

type ConnectionHandler struct {
	privateKeyPath string
	certPath       string
	tlsConfig      *tls.Config
}

func NewConnectionHandler(privateKeyPath string, certPath string) *ConnectionHandler {
	return &ConnectionHandler{privateKeyPath: privateKeyPath, certPath: certPath}
}

func (ch *ConnectionHandler) Initialize() {
	cert, err := tls.LoadX509KeyPair(ch.certPath, ch.privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	ch.tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
}

func (ch *ConnectionHandler) PerformHandshake(conn net.Conn) (net.Conn, error) {
	tlsConn := tls.Server(conn, ch.tlsConfig)
	err := tlsConn.Handshake()
	if err != nil {
		return nil, err
	}
	return tlsConn, nil
}
