package entity

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/chains-lab/profiles-svc/internal/constant"
	"github.com/chains-lab/profiles-svc/internal/dbx"
	"github.com/chains-lab/profiles-svc/internal/errx"
	"github.com/google/uuid"
)

type ProfileQ interface {
	New() dbx.ProfilesQ

	Insert(ctx context.Context, input dbx.ProfileModel) error
	Update(ctx context.Context, input map[string]any) error
	Get(ctx context.Context) (dbx.ProfileModel, error)
	Select(ctx context.Context) ([]dbx.ProfileModel, error)
	Delete(ctx context.Context) error

	FilterUserID(userID uuid.UUID) dbx.ProfilesQ
	FilterUsername(username string) dbx.ProfilesQ

	Count(ctx context.Context) (int, error)
	Page(limit, offset uint64) dbx.ProfilesQ
}

type Profiles struct {
	queries ProfileQ

	usernameRegex *regexp.Regexp
}

func NewProfile(db *sql.DB) (Profiles, error) {
	return Profiles{
		queries:       dbx.NewProfiles(db),
		usernameRegex: regexp.MustCompile(`^[a-zA-Z0-9._]{3,32}$`),
	}, nil
}

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

func (p Profiles) GenerateUsername() (string, error) {
	const (
		prefix = "citizen"
		digits = 8
	)
	buf := make([]byte, digits)
	if _, err := rand.Read(buf); err != nil {
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("cannot generate random digits: %w", err),
		)
	}
	for i := 0; i < digits; i++ {
		buf[i] = '0' + (buf[i] % 10)
	}
	return prefix + string(buf), nil
}

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

func (p Profiles) Create(ctx context.Context, userID uuid.UUID, usrnm string) error {
	_, err := p.GetByID(ctx, userID)
	if err != nil && !errors.Is(err, errx.ErrorProfileForUserDoesNotExist) {
		return err
	} else if !errors.Is(err, errx.ErrorProfileForUserDoesNotExist) {
		return errx.ErrorProfileForUserAlreadyExists.Raise(
			fmt.Errorf("profile for user '%s' already exists", userID),
		)
	}

	_, err = p.GetByUsername(ctx, usrnm)
	if !errors.Is(err, errx.ErrorProfileForUserDoesNotExist) {
		if err == nil {
			return errx.ErrorUsernameAlreadyTaken.Raise(
				fmt.Errorf("username '%s' is already taken", usrnm),
			)
		}
	}

	if err = p.ValidateUsername(usrnm); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("validating username '%s': %w", usrnm, err),
		)
	}

	createdAt := time.Now().UTC()

	err = p.queries.Insert(ctx, dbx.ProfileModel{
		UserID:    userID,
		Username:  usrnm,
		Official:  false,
		UpdatedAt: createdAt,
		CreatedAt: createdAt,
	})
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("creating profile for user '%s': %w", userID, err),
		)
	}

	return nil
}

type UpdateProfileInput struct {
	Pseudonym   *string
	Description *string
	Avatar      *string
	Official    *bool
	Sex         *string
	BirthDate   *time.Time
}

func (p Profiles) Update(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) error {
	if input.Sex != nil {
		err := constant.ParseUserSex(*input.Sex)
		if err != nil {
			return errx.ErrorSexIsNotValid.Raise(err)
		}
	}

	update := map[string]any{}

	if input.Pseudonym != nil {
		switch {
		case *input.Pseudonym == "":
			update["pseudonym"] = nil
		default:
			update["pseudonym"] = *input.Pseudonym
		}
	}

	if input.Description != nil {
		switch {
		case *input.Description == "":
			update["description"] = nil
		default:
			update["description"] = *input.Description
		}
	}

	if input.Avatar != nil {
		switch {
		case *input.Avatar == "":
			update["avatar"] = nil
		default:
			update["avatar"] = *input.Avatar
		}
	}

	if input.Sex != nil {
		switch {
		case *input.Sex == "":
			update["sex"] = nil
		default:
			update["sex"] = *input.Sex
		}
	}

	if input.BirthDate != nil {
		switch {
		case input.BirthDate.IsZero():
			update["birth_date"] = nil
		default:
			if err := p.ValidateBirthDate(*input.BirthDate); err != nil {
				return err
			}

			update["birth_date"] = *input.BirthDate
		}
	}

	if input.Official != nil {
		update["official"] = *input.Official
	}

	err := p.queries.FilterUserID(userID).Update(ctx, update)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorProfileForUserDoesNotExist.Raise(
				fmt.Errorf("profile for user '%s' does not exist", userID),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating profile for user '%s': %w", userID, err),
			)
		}
	}

	return nil
}

func (p Profiles) UpdateUsername(ctx context.Context, userID uuid.UUID, usrnm string) error {
	if err := p.ValidateUsername(usrnm); err != nil {
		return errx.ErrorUsernameIsNotValid.Raise(
			fmt.Errorf("validating username '%s': %w", usrnm, err),
		)
	}

	now := time.Now().UTC()

	profile, err := p.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = p.GetByUsername(ctx, usrnm)
	if !errors.Is(err, errx.ErrorProfileForUserDoesNotExist) {
		return err
	}
	if err == nil {
		return errx.ErrorUsernameAlreadyTaken.Raise(
			fmt.Errorf("username '%s' is already taken", usrnm),
		)
	}

	if profile.Username == usrnm {
		return nil // No change needed
	}

	err = p.queries.FilterUserID(userID).Update(ctx, map[string]any{
		"username":   usrnm,
		"updated_at": now,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorProfileForUserDoesNotExist.Raise(
				fmt.Errorf("profile for user '%s' does not exist", userID),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating username for user '%s': %w", userID, err),
			)
		}
	}

	return nil
}

func (p Profiles) GetByID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	profile, err := p.queries.FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Profile{}, errx.ErrorProfileForUserDoesNotExist.Raise(
				fmt.Errorf("profile for user '%s' does not exist", userID),
			)
		default:
			return models.Profile{}, errx.ErrorInternal.Raise(
				fmt.Errorf("getting profile for user '%s': %w", userID, err),
			)
		}
	}

	return ProfileFromDb(profile), nil
}

func (p Profiles) GetByUsername(ctx context.Context, username string) (models.Profile, error) {
	profile, err := p.queries.FilterUsername(username).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Profile{}, errx.ErrorProfileForUserDoesNotExist.Raise(
				fmt.Errorf("profile with username '%s' does not exist", username),
			)
		default:
			return models.Profile{}, errx.ErrorInternal.Raise(
				fmt.Errorf("getting profile with username '%s': %w", username, err),
			)
		}
	}

	return ProfileFromDb(profile), nil
}

func ProfileFromDb(input dbx.ProfileModel) models.Profile {
	return models.Profile{
		UserID:      input.UserID,
		Username:    input.Username,
		Pseudonym:   input.Pseudonym,
		Description: input.Description,
		Avatar:      input.Avatar,
		Official:    input.Official,
		Sex:         input.Sex,
		Birthdate:   input.BirthDate,
		UpdatedAt:   input.UpdatedAt,
		CreatedAt:   input.CreatedAt,
	}
}
