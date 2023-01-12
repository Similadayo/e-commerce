package utils

import (
	"errors"
	"regexp"

	"github.com/Similadayo/backend/models"
)

func Validate(user *models.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	if user.Address == "" {
		return errors.New("address is required")
	}
	if user.PhoneNumber == "" {
		return errors.New("phone number is required")
	}
	return nil
}

// IsEmailValid check if the given string is a valid email
func IsEmailValid(email string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

// IsPasswordValid check if the given string is a valid password
func IsPasswordValid(password string) bool {
	return len(password) >= 6
}
