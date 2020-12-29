package models

import (
	"errors"
	"fmt"
	"time"
)

var (
	//ErrConnNotInit ...
	ErrConnNotInit = errors.New("connection does not init")
)

// Prefix ...
func (c *Connection) Prefix() error {

	if c.Test {
		return nil
	}

	if c.Conn == nil {
		return ErrConnNotInit
	}

	date := time.Now().Format("2006-01-02 15:04:05")

	_, err := (*c.Conn).Write([]byte(fmt.Sprintf("[%v][%v]:", date, c.Name)))

	return err
}

// Write ...
func (c *Connection) Write(s string) error {

	if c.Test {
		return nil
	}

	if c.Conn == nil {
		return ErrConnNotInit
	}

	_, err := (*c.Conn).Write([]byte(s + "\n"))

	return err
}
