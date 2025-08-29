package entity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/chains-lab/profiles-svc/internal/app/username"
	"github.com/chains-lab/profiles-svc/internal/dbx"
	"github.com/chains-lab/profiles-svc/internal/errx"
	"github.com/google/uuid"
)

type ProfileQ interface {
	New() dbx.ProfilesQ

	Insert(ctx context.Context, input dbx.ProfileModel) error
	Update(ctx context.Context, input dbx.UpdateProfileInput) error
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
}

func NewProfile(db *sql.DB) (Profiles, error) {
	return Profiles{
		queries: dbx.NewProfiles(db),
	}, nil
}

type CreateProfileInput struct {
	Username    string
	Pseudonym   *string
	Description *string
	Avatar      *string
	Sex         *string
	BirthDate   *time.Time
}

func (p Profiles) Create(ctx context.Context, userID uuid.UUID, input CreateProfileInput) error {
	_, err := p.GetByID(ctx, userID)
	if !errors.Is(err, errx.ErrorProfileForUserDoesNotExist) {
		return err
	}

	_, err = p.GetByUsername(ctx, input.Username)
	if !errors.Is(err, errx.ErrorProfileForUserDoesNotExist) {
		if err == nil {
			return errx.ErrorUsernameAlreadyTaken.Raise(
				fmt.Errorf("username '%s' is already taken", input.Username),
			)
		}
	}

	if err = username.ValidateUsername(input.Username); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("validating username '%s': %w", input.Username, err),
		)
	}

	createdAt := time.Now().UTC()

	err = p.queries.Insert(ctx, dbx.ProfileModel{
		UserID:            userID,
		Username:          input.Username,
		Pseudonym:         input.Pseudonym,
		Description:       input.Description,
		Avatar:            input.Avatar,
		Official:          false,
		Sex:               input.Sex,
		BirthDate:         input.BirthDate,
		UsernameUpdatedAt: createdAt,
		UpdatedAt:         createdAt,
		CreatedAt:         createdAt,
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
		//TODO: Validate sex value
		//if err != nil {
		//	return errx.RaiseSexIsNotValid(err)
		//}
	}

	err := p.queries.FilterUserID(userID).Update(ctx, dbx.UpdateProfileInput{
		Pseudonym:   input.Pseudonym,
		Description: input.Description,
		Avatar:      input.Avatar,
		Official:    input.Official,
		Sex:         input.Sex,
		BirthDate:   input.BirthDate,
		UpdatedAt:   time.Now().UTC(),
	})
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

func (p Profiles) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {
	if err := username.ValidateUsername(username); err != nil {
		return errx.ErrorUsernameIsNotValid.Raise(
			fmt.Errorf("validating username '%s': %w", username, err),
		)
	}

	now := time.Now().UTC()

	profile, err := p.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	elapsed := now.Sub(profile.UsernameUpdatedAt)
	if elapsed < 14*24*time.Hour {
		return errx.ErrorUsernameUpdateCooldown.Raise(
			fmt.Errorf("username was updated %.0f hours ago", elapsed.Hours()),
		)
	}

	_, err = p.GetByUsername(ctx, username)
	if !errors.Is(err, errx.ErrorProfileForUserDoesNotExist) {
		return err
	}
	if err == nil {
		return errx.ErrorUsernameAlreadyTaken.Raise(
			fmt.Errorf("username '%s' is already taken", username),
		)
	}

	if profile.Username == username {
		return nil // No change needed
	}

	err = p.queries.FilterUserID(userID).Update(ctx, dbx.UpdateProfileInput{
		Username:          &username,
		UsernameUpdatedAt: &now,
		UpdatedAt:         now,
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
		UserID:            input.UserID,
		Username:          input.Username,
		Pseudonym:         input.Pseudonym,
		Description:       input.Description,
		Avatar:            input.Avatar,
		Official:          input.Official,
		Sex:               input.Sex,
		BirthDate:         input.BirthDate,
		UsernameUpdatedAt: input.UsernameUpdatedAt,
		UpdatedAt:         input.UpdatedAt,
		CreatedAt:         input.CreatedAt,
	}
}
