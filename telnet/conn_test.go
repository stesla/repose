package telnet_test

import (
	"bytes"
	"errors"
	"net"
	"repose/telnet"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestRead(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		buf      string
		expected string
	}{
		{"foo", "foo"},
		{"foo\xf1bar", "foo\xf1bar"},
		{"foo\xff\xf1bar", "foobar"},
		{"foo\xff\xffbar", "foo\xffbar"},
		{"foo\r\x00bar", "foo\rbar"},
		{"foo\r\nbar", "foo\nbar"},
		{"foo\rbar", "foobar"},
	}
	for _, test := range tests {
		var conn net.Conn = telnet.NewConn(&conn{
			r: bytes.NewBufferString(test.buf),
		})
		buf := make([]byte, len(test.buf))
		n, err := conn.Read(buf)
		assert.NoError(err)
		assert.Equal(len(test.expected), n)
		assert.Equal(test.expected, string(buf[:n]))
	}
}

func TestWrite(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		buf      string
		expected string
	}{
		{"foo", "foo"},
		{"foo\xf1bar", "foo\xf1bar"},
		{"foo\xffbar", "foo\xff\xffbar"},
		{"foo\nbar", "foo\r\nbar"},
		{"foo\rbar", "foo\r\x00bar"},
	}
	for _, test := range tests {
		var buf bytes.Buffer
		var conn net.Conn = telnet.NewConn(&conn{w: &buf})
		n, err := conn.Write([]byte(test.buf))
		assert.NoError(err)
		assert.Equal(len(test.buf), n)
		assert.Equal(test.expected, buf.String())
	}
}

type writeErrorTest struct {
	errorAfter int
	expected   int
}

func (w writeErrorTest) Write(p []byte) (n int, err error) {
	return w.errorAfter, errors.New("failed to write")
}

func TestWriteError(t *testing.T) {
	assert := assert.New(t)
	buf := "1\n2\r3\xff4"
	tests := []writeErrorTest{
		{0, 0},  // nothing
		{1, 1},  // 1
		{2, 1},  // CR
		{3, 2},  // LF
		{4, 3},  // 2
		{5, 3},  // CR
		{6, 4},  // NUL
		{7, 5},  // 3
		{8, 5},  // IAC
		{9, 6},  // IAC
		{10, 7}, // 4
	}
	for _, test := range tests {
		var conn net.Conn = telnet.NewConn(&conn{w: test})
		n, err := conn.Write([]byte(buf))
		assert.Error(err)
		assert.Equal(test.expected, n)
	}
}
