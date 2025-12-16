package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/google/uuid"
)

type UpdateProfileParams struct {
	Pseudonym   *string
	Description *string
	Avatar      *string
}

func (s Service) UpdateProfile(ctx context.Context, accountID uuid.UUID, input UpdateProfileParams) (entity.Profile, error) {
	p, err := s.GetProfileByID(ctx, accountID)
	if err != nil {
		return entity.Profile{}, err
	}

	if input == (UpdateProfileParams{}) {
		return p, nil
	}

	profile, err := s.db.UpdateProfile(ctx, accountID, input)
	if err != nil {
		return entity.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("updating profile for user '%s': %w", accountID, err),
		)
	}

	return profile, nil
}

func (s Service) UpdateProfileOfficial(ctx context.Context, accountID uuid.UUID, official bool) (entity.Profile, error) {
	_, err := s.GetProfileByID(ctx, accountID)
	if err != nil {
		return entity.Profile{}, err
	}

	profile, err := s.db.UpdateProfileOfficial(ctx, accountID, official)
	if err != nil {
		return entity.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("updating profile for user '%s': %w", accountID, err),
		)
	}

	return profile, nil
}

func (s Service) UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (entity.Profile, error) {
	profile, err := s.db.UpdateProfileUsername(ctx, accountID, username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return entity.Profile{}, errx.ErrorProfileNotFound.Raise(
				fmt.Errorf("profile for user '%s' does not exist", accountID),
			)
		default:
			return entity.Profile{}, errx.ErrorInternal.Raise(
				fmt.Errorf("updating username for user '%s': %w", accountID, err),
			)
		}
	}

	return profile, nil
}
