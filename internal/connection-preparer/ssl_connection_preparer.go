package connection_preparer

import (
	"crypto/tls"
	"log"
	"net"
)

type SslConnectionPreparer struct {
	privateKeyPath string
	certPath       string
	tlsConfig      *tls.Config
}

func NewConnectionHandler(privateKeyPath string, certPath string) *SslConnectionPreparer {
	return &SslConnectionPreparer{privateKeyPath: privateKeyPath, certPath: certPath}
}

func (scp *SslConnectionPreparer) Initialize() {
	cert, err := tls.LoadX509KeyPair(scp.certPath, scp.privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	scp.tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
}

func (scp *SslConnectionPreparer) Prepare(conn net.Conn) (net.Conn, error) {
	tlsConn := tls.Server(conn, scp.tlsConfig)
	err := tlsConn.Handshake()
	if err != nil {
		return nil, err
	}
	return tlsConn, nil
}
