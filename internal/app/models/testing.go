package models

import (
	"errors"
	"kitKat/internal/app/utils"
	"net"

	// "netcat/internal/app/utils"
	"testing"
	"time"
)

// TestConnection ...
func TestConnection(t *testing.T, str, port string, rooms int) (*Connection, error) {

	if str == "" || port == "" || rooms <= 0 {
		return nil, errors.New("invalid input")
	}

	c, err := TestConn(t, port)

	if err != nil {
		return nil, err
	}

	if !utils.ValidName(str) {
		return nil, errors.New("invalid name")
	}

	return &Connection{
		Name:    str,
		Conn:    c,
		Action:  time.Now(),
		Room:    1,
		LastMsg: make([]string, rooms+1),
	}, nil
}

// TestConn ...
func TestConn(t *testing.T, port string) (*net.Conn, error) {

	c, err := net.Dial("tcp", port)

	return &c, err
}

// NewTestConn ...
func NewTestConn(t *testing.T, name string, rooms int) (*Connection, error) {

	if name == "" || rooms <= 0 {
		return nil, errors.New("invalid input")
	}

	if !utils.ValidName(name) {
		return nil, errors.New("invalid name")
	}

	c := TestNetConn(t)

	return &Connection{
		Name:    name,
		Conn:    &c,
		Action:  time.Now(),
		Room:    1,
		LastMsg: make([]string, rooms+1),
	}, nil
}
