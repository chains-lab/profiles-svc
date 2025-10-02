package profile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/chains-lab/profiles-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Create(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error) {
	_, err := s.GetByID(ctx, userID)
	if err != nil && !errors.Is(err, errx.ErrorProfileNotFound) {
		return models.Profile{}, err
	} else if !errors.Is(err, errx.ErrorProfileNotFound) {
		return models.Profile{}, errx.ErrorProfileAlreadyExists.Raise(
			fmt.Errorf("profile for user '%s' already exists", userID),
		)
	}

	_, err = s.GetByUsername(ctx, username)
	if !errors.Is(err, errx.ErrorProfileNotFound) {
		if err == nil {
			return models.Profile{}, errx.ErrorUsernameAlreadyTaken.Raise(
				fmt.Errorf("username '%s' is already taken", username),
			)
		}
	}

	if err = validateUsername(username); err != nil {
		return models.Profile{}, errx.ErrorUsernameIsNotValid.Raise(
			fmt.Errorf("validating username '%s': %w", username, err),
		)
	}

	createdAt := time.Now().UTC()

	p := models.Profile{
		UserID:    userID,
		Username:  username,
		Official:  false,
		UpdatedAt: createdAt,
		CreatedAt: createdAt,
	}

	err = s.db.CreateProfile(ctx, p)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("creating profile for user '%s': %w", userID, err),
		)
	}

	return p, nil
}
