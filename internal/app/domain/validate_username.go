package domain

import (
	"errors"
	"regexp"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._]{3,32}$`)
)

func ValidateUsername(username string) error {
	length := len(username)
	if length < 3 || length > 32 {
		return errors.New("username must be between 3 and 32 characters")
	}

	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain Latin letters, digits, dot (.) and underscore (_)")
	}

	return nil
}
