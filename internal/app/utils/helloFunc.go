package utils

import (
	"net"
	"os"
)

// Hello func prints welcome message for each connection
func Hello(c *net.Conn) error {

	file, err := os.Open("./static/pinguin.txt")
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}
	size := info.Size()
	buffer := make([]byte, size)
	if _, err := file.Read(buffer); err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	if _, err := (*c).Write([]byte("Welcome to TCP-Chat!\n")); err != nil {
		return err
	}

	if _, err := (*c).Write(buffer); err != nil {
		return err
	}

	if _, err := (*c).Write([]byte("\n[ENTER YOUR NAME]: ")); err != nil {
		return err
	}
	return nil
}
