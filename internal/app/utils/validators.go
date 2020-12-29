package utils

import "regexp"

// ValidName func checks name of user
func ValidName(s string) bool {
	re := regexp.MustCompile(`^\w{1,10}$`)
	return re.MatchString(s)
}
