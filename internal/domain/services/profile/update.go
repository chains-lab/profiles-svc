package profile

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/profiles-svc/internal/domain/enum"
	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/chains-lab/profiles-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Update struct {
	Pseudonym   *string
	Description *string
	Avatar      *string
	Sex         *string
	BirthDate   *time.Time
}

func (s Service) Update(ctx context.Context, userID uuid.UUID, input Update) (models.Profile, error) {
	p, err := s.GetByID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	if input == (Update{}) {
		return p, nil
	}

	now := time.Now().UTC()

	if input.Pseudonym != nil {
		if *input.Pseudonym == "" {
			p.Pseudonym = nil
		} else {
			p.Pseudonym = input.Pseudonym
		}
	}
	if input.Description != nil {
		if *input.Description == "" {
			p.Description = nil
		} else {
			p.Description = input.Description
		}
	}
	if input.Avatar != nil {
		if *input.Avatar == "" {
			p.Avatar = nil
		} else {
			p.Avatar = input.Avatar
		}
	}
	if input.Sex != nil {
		if err = enum.CheckUserSexes(*input.Sex); err != nil {
			return models.Profile{}, errx.ErrorSexIsNotValid.Raise(err)
		}
		p.Sex = input.Sex
	}
	if input.BirthDate != nil {
		if err = validateBirthDate(*input.BirthDate); err != nil {
			return models.Profile{}, err
		}
		p.BirthDate = input.BirthDate
	}

	p.UpdatedAt = now

	err = s.db.UpdateProfile(ctx, userID, input, now)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("updating profile for user '%s': %w", userID, err),
		)
	}

	return p, nil
}

func (s Service) UpdateOfficial(ctx context.Context, userID uuid.UUID, official bool) (models.Profile, error) {
	p, err := s.GetByID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	now := time.Now().UTC()

	err = s.db.UpdateProfileOfficial(ctx, userID, official, now)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("updating official profile for user '%s': %w", userID, err),
		)
	}

	p.Official = official
	p.UpdatedAt = now

	return p, nil
}

func (s Service) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error) {
	p, err := s.GetByID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	if err = validateUsername(username); err != nil {
		return models.Profile{}, errx.ErrorUsernameIsNotValid.Raise(
			fmt.Errorf("validating username '%s': %w", username, err),
		)
	}

	now := time.Now().UTC()

	_, err = s.GetByUsername(ctx, username)
	if !errors.Is(err, errx.ErrorProfileNotFound) {
		return models.Profile{}, err
	}
	if err == nil {
		return models.Profile{}, errx.ErrorUsernameAlreadyTaken.Raise(
			fmt.Errorf("username '%s' is already taken", username),
		)
	}

	if p.Username == username {
		return models.Profile{}, nil
	}

	err = s.db.UpdateProfileUsername(ctx, userID, username, now)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Profile{}, errx.ErrorProfileNotFound.Raise(
				fmt.Errorf("profile for user '%s' does not exist", userID),
			)
		default:
			return models.Profile{}, errx.ErrorInternal.Raise(
				fmt.Errorf("updating username for user '%s': %w", userID, err),
			)
		}
	}

	p.UpdatedAt = now
	p.Username = username

	return p, nil
}
