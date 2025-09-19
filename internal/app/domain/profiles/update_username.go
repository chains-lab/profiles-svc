package profiles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/profiles-svc/internal/errx"
	"github.com/google/uuid"
)

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
