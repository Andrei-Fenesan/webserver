package connection_preparer

import "net"

type PlainConnectionPreparer struct{}

func (p *PlainConnectionPreparer) Prepare(conn net.Conn) (net.Conn, error) {
	return conn, nil
}
