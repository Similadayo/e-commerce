package utils

import (
	"regexp"
)

// IsEmailValid check if the given string is a valid email
func IsEmailValid(email string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

// IsPasswordValid check if the given string is a valid password
func IsPasswordValid(password string) bool {
	return len(password) >= 6
}
