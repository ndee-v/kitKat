package models

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

// Addr ... for testing
type Addr interface {
	Network() string // name of the network (for example, "tcp", "udp")
	String() string  // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
}

// MyAddr ...testing
type MyAddr struct {
}

// Network ...testing
func (a MyAddr) Network() string {
	return ""
}

// String ...testing
func (a MyAddr) String() string {
	return ""
}

// NetConnection using for test cases
type NetConnection struct {
	ReadLimit  bool
	ReadsLeft  int
	WriteLimit bool
	WritesLeft int
	Open       bool
	Buff       []byte
}

//NetConnection ...

// TestNetConn ...
func TestNetConn(t *testing.T) net.Conn {
	return &NetConnection{
		Open: true,
	}
}

//TestNetConnWriteLimited ...
func TestNetConnWriteLimited(t *testing.T, num int) (net.Conn, error) {

	if num <= 0 {
		return nil, errors.New("invalid number of limit")
	}

	return &NetConnection{
		Open:       true,
		WriteLimit: true,
		WritesLeft: num,
	}, nil
}

// Read ...
func (n *NetConnection) Read(b []byte) (int, error) {

	if !n.Open {
		return 0, errors.New("connection closed")
	}

	for i := range n.Buff {
		if i > len(b)-1 {
			return i, io.ErrUnexpectedEOF
		}
		b[i] = n.Buff[i]
	}

	l := len(n.Buff)

	n.Buff = nil

	return l, nil

}

// Write ...
func (n *NetConnection) Write(b []byte) (int, error) {

	if n.WriteLimit {
		if n.WritesLeft == 0 {
			n.Open = false
			return 0, errors.New("connection closed")
		}
		n.WritesLeft--
	}

	if !n.Open {
		return 0, errors.New("connection closed")
	}

	startLen := len(n.Buff)
	n.Buff = append(n.Buff, b...)
	writedLen := len(n.Buff) - startLen

	if writedLen != len(b) {
		return writedLen, io.ErrUnexpectedEOF
	}

	return writedLen, nil
}

// Close ...
func (n *NetConnection) Close() error {

	if !n.Open {
		return errors.New("connection closed")
	}

	n.Open = false

	return nil
}

// LocalAddr ...
func (n *NetConnection) LocalAddr() net.Addr {
	a := MyAddr{}
	return a

}

// RemoteAddr ...
func (n *NetConnection) RemoteAddr() net.Addr {
	a := MyAddr{}
	return a
}

// SetDeadline ...
func (n *NetConnection) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline ...
func (n *NetConnection) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline ...
func (n *NetConnection) SetWriteDeadline(t time.Time) error {
	return nil
}
