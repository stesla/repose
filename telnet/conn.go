package telnet

import "net"

const (
	bIAC = '\xff'
)

type TelnetConnection struct {
	net.Conn
}

func NewConn(conn net.Conn) *TelnetConnection {
	return &TelnetConnection{conn}
}

func (t *TelnetConnection) Read(p []byte) (int, error) {
	b := make([]byte, len(p))
	bn, err := t.Conn.Read(b)
	var i, j int
	for i, j = 0, 0; i < bn && j < len(p); i++ {
		switch b[i] {
		case bIAC:
			i++
			if i < bn && b[i] == bIAC {
				p[j] = bIAC
				j++
			}
		default:
			p[j] = b[i]
			j++
		}
	}
	return j, err
}
