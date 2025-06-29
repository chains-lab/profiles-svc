package entities

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
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

func (p Profiles) Create(ctx context.Context, userID uuid.UUID) error {
	createdAt := time.Now().UTC()

	username, err := func() (string, error) {
		const suffixLen = 12
		max := new(big.Int).Exp(big.NewInt(10), big.NewInt(suffixLen), nil)
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("генерация случайного суффикса для username: %w", err)
		}
		return "elector" + fmt.Sprintf("%0*d", suffixLen, n), nil
	}()
	if err != nil {
		return err
	}

	err = p.queries.Insert(ctx, dbx.ProfileModel{
		UserID:    userID,
		Username:  username,
		Official:  false,
		UpdatedAt: createdAt,
		CreatedAt: createdAt,
	})
	if err != nil {
		return ape.ErrorInternal(err)
	}

	return nil
}

type UpdateProfileInput struct {
	Username    *string `json:"username,omitempty"`
	Pseudonym   *string `json:"pseudonym,omitempty"`
	Description *string `json:"description,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
	Official    *bool   `json:"official,omitempty"`
}

func (p Profiles) Update(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) error {
	if input.Username != nil {
		_, err := p.queries.New().FilterUsername(*input.Username).Get(ctx)
		if errors.Is(err, sql.ErrNoRows) {

		} else if err != nil {
			return ape.ErrorInternal(err)
		}

		if err == nil {
			return ape.ErrorUsernameAlreadyTaken(err, *input.Username)
		}
	}

	err := p.queries.FilterUserID(userID).Update(ctx, dbx.UpdateProfileInput{
		Username:    input.Username,
		Pseudonym:   input.Pseudonym,
		Description: input.Description,
		Avatar:      input.Avatar,
		Official:    input.Official,
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		}
		return ape.ErrorInternal(err)
	}

	return nil
}

func (p Profiles) GetByID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	profile, err := p.queries.FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Profile{}, ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		default:
			return models.Profile{}, ape.ErrorInternal(err)
		}
	}

	return ProfileFromDb(profile), nil
}

func (p Profiles) GetByUsername(ctx context.Context, username string) (models.Profile, error) {
	profile, err := p.queries.FilterUsername(username).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Profile{}, ape.ErrorCabinetForUserDoesNotExist(err, username)
		default:
			return models.Profile{}, ape.ErrorInternal(err)
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
