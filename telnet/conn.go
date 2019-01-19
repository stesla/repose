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

func (t *TelnetConnection) Read(p []byte) (i int, err error) {
	b := make([]byte, len(p))
	n, err := t.Conn.Read(b)
	b = b[:n]
	var f parseState = parseDefault
	for i = 0; i < len(p) && len(b) > 0; b = b[1:] {
		i, f = f(b[0], p, i)
	}
	return
}

type parseState func(byte, []byte, int) (int, parseState)

func parseDefault(c byte, p []byte, j int) (int, parseState) {
	switch c {
	case bIAC:
		return j, parseIAC
	case '\r':
		return j, parseCR
	default:
		p[j] = c
		return j + 1, parseDefault
	}
}

func parseIAC(c byte, p []byte, j int) (int, parseState) {
	switch c {
	case bIAC:
		p[j] = c
		j++
	}
	return j, parseDefault
}

func parseCR(c byte, p []byte, j int) (int, parseState) {
	if c == '\x00' {
		p[j] = '\r'
	} else {
		p[j] = c
	}
	return j + 1, parseDefault
}

func (t *TelnetConnection) Write(p []byte) (n int, err error) {
	b := make([]byte, 0, 2*len(p))
	nmap := map[int]int{}
	for i, c := range p {
		n := len(b)
		switch c {
		case bIAC:
			b = append(b, bIAC, bIAC)
		case '\n':
			b = append(b, '\r', '\n')
		case '\r':
			b = append(b, '\r', '\x00')
		default:
			b = append(b, c)
		}
		switch len(b) - n {
		case 2:
			nmap[len(b)-1] = i
			fallthrough
		case 1:
			nmap[len(b)] = i + 1
		}
	}
	n, err = t.Conn.Write(b)
	n = nmap[n]
	return
}
