package profiles

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/profiles-svc/internal/dbx"
	"github.com/chains-lab/profiles-svc/internal/errx"
	"github.com/google/uuid"
)

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
