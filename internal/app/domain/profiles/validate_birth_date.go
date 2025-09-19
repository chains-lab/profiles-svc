package profiles

import (
	"fmt"
	"time"

	"github.com/chains-lab/profiles-svc/internal/errx"
)

func (p Profiles) ValidateBirthDate(birthDate time.Time) error {
	minDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	if birthDate.Before(minDate) {
		return errx.ErrorBirthdateIsNotValid.Raise(
			fmt.Errorf("birthdate date must be in the future (1900-01-01)"),
		)
	}

	now := time.Now().UTC()
	age := now.Year() - birthDate.Year()
	if now.Month() < birthDate.Month() ||
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	if age < 12 {
		return errx.ErrorUserTooYoung.Raise(
			fmt.Errorf("user must be at least 12 years old"),
		)
	}

	return nil
}
