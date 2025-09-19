package profiles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/chains-lab/profiles-svc/internal/errx"
	"github.com/google/uuid"
)

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
