package profile

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/chains-lab/profiles-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Service struct {
	db database
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._]{3,32}$`)

func validateUsername(username string) error {
	length := len(username)
	if length < 3 || length > 32 {
		return fmt.Errorf("username must be between 3 and 32 characters")
	}

	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("username can only contain Latin letters, digits, dot (.) and underscore (_)")
	}

	return nil
}

func validateBirthDate(birthDate time.Time) error {
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

func New(db database) Service {
	return Service{
		db: db,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreateProfile(ctx context.Context, profile models.Profile) error

	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)

	FilterProfiles(ctx context.Context, filters FilterParams, page uint64, size uint64) (models.ProfileCollection, error)

	UpdateProfile(ctx context.Context, userID uuid.UUID, input Update, updatedAt time.Time) error
	UpdateProfileUsername(ctx context.Context, userID uuid.UUID, username string, updatedAt time.Time) error
	UpdateProfileBirthDate(ctx context.Context, userID uuid.UUID, birthDate time.Time, updatedAt time.Time) error
	UpdateProfileSex(ctx context.Context, userID uuid.UUID, sex string, updatedAt time.Time) error
	UpdateProfileOfficial(ctx context.Context, userID uuid.UUID, official bool, updatedAt time.Time) error
}
