package utils

import (
	"fmt"
	"time"
)

// PreText func prints date and author of every message
func PreText(s string) string {
	date := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%v][%v]:", date, s)
	//return "[" + time.Now().Format("2006-01-02 15:04:05") + "][" + s + "]:"
}

// Prefix func prints date and author of every message
func Prefix(s string) string {
	date := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%v][%v]:", date, s)
}
