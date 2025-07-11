package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/ape"
	"github.com/chains-lab/elector-cab-svc/internal/app/domain"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/elector-cab-svc/internal/dbx"
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
}

func (p Profiles) Create(ctx context.Context, userID uuid.UUID, input CreateProfileInput) error {
	_, err := p.GetByID(ctx, userID)
	if !errors.Is(err, ape.ErrorProfileForUserDoesNotExist) {
		return err
	}

	_, err = p.GetByUsername(ctx, input.Username)
	if !errors.Is(err, ape.ErrorProfileForUserDoesNotExist) {
		if err == nil {
			return ape.RaiseUsernameAlreadyTaken(err, input.Username)
		}
	}

	if err := domain.ValidateUsername(input.Username); err != nil {
		return ape.RaiseInternal(err)
	}

	createdAt := time.Now().UTC()

	err = p.queries.Insert(ctx, dbx.ProfileModel{
		UserID:            userID,
		Username:          input.Username,
		Pseudonym:         input.Pseudonym,
		Description:       input.Description,
		Avatar:            input.Avatar,
		Official:          false,
		UsernameUpdatedAt: createdAt,
		UpdatedAt:         createdAt,
		CreatedAt:         createdAt,
	})
	if err != nil {
		return ape.RaiseInternal(err)
	}

	return nil
}

type UpdateProfileInput struct {
	Pseudonym   *string `json:"pseudonym,omitempty"`
	Description *string `json:"description,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
	Official    *bool   `json:"official,omitempty"`
}

func (p Profiles) Update(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) error {

	err := p.queries.FilterUserID(userID).Update(ctx, dbx.UpdateProfileInput{
		Pseudonym:   input.Pseudonym,
		Description: input.Description,
		Avatar:      input.Avatar,
		Official:    input.Official,
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseProfileForUserNotFound(err, userID)
		}
		return ape.RaiseInternal(err)
	}

	return nil
}

func (p Profiles) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {
	if err := domain.ValidateUsername(username); err != nil {
		return ape.RaiseUsernameIsNotValid(err)
	}

	now := time.Now().UTC()

	profile, err := p.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	elapsed := now.Sub(profile.UsernameUpdatedAt)
	if elapsed < 14*24*time.Hour {
		return ape.RaiseUsernameUpdateCooldown(
			fmt.Errorf("username was updated %.0f hours ago", elapsed.Hours()),
			profile.UserID,
		)
	}

	_, err = p.GetByUsername(ctx, username)
	if !errors.Is(err, ape.ErrorProfileForUserDoesNotExist) {
		return err
	}
	if err == nil {
		return ape.RaiseUsernameAlreadyTaken(err, username)
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
			return ape.RaiseProfileForUserNotFound(err, userID)
		default:
			return ape.RaiseInternal(err)
		}
	}

	return nil
}

func (p Profiles) GetByID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	profile, err := p.queries.FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Profile{}, ape.RaiseProfileForUserNotFound(err, userID)
		default:
			return models.Profile{}, ape.RaiseInternal(err)
		}
	}

	return ProfileFromDb(profile), nil
}

func (p Profiles) GetByUsername(ctx context.Context, username string) (models.Profile, error) {
	profile, err := p.queries.FilterUsername(username).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Profile{}, ape.RaiseProfileForUserNotFoundByUsername(err, username)
		default:
			return models.Profile{}, ape.RaiseInternal(err)
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
		UpdatedAt:   input.UpdatedAt,
		CreatedAt:   input.CreatedAt,
	}
}
