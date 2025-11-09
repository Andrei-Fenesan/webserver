package connection_preparer

import "net"

type ConnectionPreparer interface {
	Prepare(conn net.Conn) (net.Conn, error)
}
