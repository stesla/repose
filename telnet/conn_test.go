package telnet_test

import (
	"bytes"
	"net"
	"repose/telnet"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestRead(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		input    string
		expected string
	}{
		{"foo", "foo"},
		{"foo\xff\xf1bar", "foobar"},
		{"foo\xff\xffbar", "foo\xffbar"},
	}
	for _, test := range tests {
		var conn net.Conn = telnet.NewConn(&conn{
			rbuf: bytes.NewBufferString(test.input),
		})
		buf := make([]byte, len(test.input))
		n, err := conn.Read(buf)
		assert.NoError(err)
		assert.Equal(len(test.expected), n)
		assert.Equal(test.expected, string(buf[:n]))
	}
}
