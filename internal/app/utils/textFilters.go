package utils

import (
	"errors"
	"regexp"
)

// TextFilter func prepare passed text to normal, readable format
func TextFilter(s string) (string, error) {
	s = s[:len(s)-1]
	// one or whitespaces at start of string of at end of string
	re := regexp.MustCompile(`(^\s+)|(\s+$)`)
	//one or more whitespaces into string
	re2 := regexp.MustCompile(`\s+`)
	// delete groups of whitespaces
	s = re.ReplaceAllString(s, "")
	s = re2.ReplaceAllString(s, " ")
	if len(s) == 0 {
		return "", errors.New("invalid input")
	}
	return s, nil
}
