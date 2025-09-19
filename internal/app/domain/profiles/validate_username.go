package profiles

import "errors"

func (p Profiles) ValidateUsername(username string) error {
	length := len(username)
	if length < 3 || length > 32 {
		return errors.New("username must be between 3 and 32 characters")
	}

	if !p.usernameRegex.MatchString(username) {
		return errors.New("username can only contain Latin letters, digits, dot (.) and underscore (_)")
	}

	return nil
}
