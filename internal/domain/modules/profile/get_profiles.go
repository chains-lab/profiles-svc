package profile

import (
	"context"
	"fmt"

	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) GetProfileByID(ctx context.Context, userID uuid.UUID) (entity.Profile, error) {
	profile, err := s.db.GetProfileByAccountID(ctx, userID)
	if err != nil {
		return entity.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting profile for user '%s': %w", userID, err),
		)
	}

	if profile.IsNil() {
		return entity.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("profile for user '%s' does not exist", userID),
		)
	}

	return profile, nil
}

func (s Service) GetProfileByUsername(ctx context.Context, username string) (entity.Profile, error) {
	profile, err := s.db.GetProfileByUsername(ctx, username)
	if err != nil {
		return entity.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting profile with username '%s': %w", username, err),
		)
	}

	if profile.IsNil() {
		return entity.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("profile with username '%s' does not exist", username),
		)
	}

	return profile, nil
}
