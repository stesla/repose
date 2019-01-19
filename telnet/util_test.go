package telnet_test

import (
	"io"
	"net"
	"time"
)

type conn struct {
	r io.Reader
	w io.Writer
}

func (c *conn) Read(b []byte) (n int, err error) {
	n, err = c.r.Read(b)
	return
}

func (c *conn) Write(b []byte) (n int, err error) {
	n, err = c.w.Write(b)
	return
}

func (c *conn) Close() error {
	panic("not implemented")
}

func (c *conn) LocalAddr() net.Addr {
	panic("not implemented")
}

func (c *conn) RemoteAddr() net.Addr {
	panic("not implemented")
}

func (c *conn) SetDeadline(t time.Time) error {
	panic("not implemented")
}

func (c *conn) SetReadDeadline(t time.Time) error {
	panic("not implemented")
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	panic("not implemented")
}
