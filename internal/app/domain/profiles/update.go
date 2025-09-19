package profiles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/profiles-svc/internal/errx"
	"github.com/google/uuid"
)

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
		err := enum.CheckUserSexes(*input.Sex)
		if err != nil {
			return errx.ErrorSexIsNotValid.Raise(err)
		}
	}

	update := map[string]any{}

	if input.Pseudonym != nil {
		switch {
		case *input.Pseudonym == "":
			update["pseudonym"] = nil
		default:
			update["pseudonym"] = *input.Pseudonym
		}
	}

	if input.Description != nil {
		switch {
		case *input.Description == "":
			update["description"] = nil
		default:
			update["description"] = *input.Description
		}
	}

	if input.Avatar != nil {
		switch {
		case *input.Avatar == "":
			update["avatar"] = nil
		default:
			update["avatar"] = *input.Avatar
		}
	}

	if input.Sex != nil {
		switch {
		case *input.Sex == "":
			update["sex"] = nil
		default:
			update["sex"] = *input.Sex
		}
	}

	if input.BirthDate != nil {
		switch {
		case input.BirthDate.IsZero():
			update["birth_date"] = nil
		default:
			if err := p.ValidateBirthDate(*input.BirthDate); err != nil {
				return err
			}

			update["birth_date"] = *input.BirthDate
		}
	}

	if input.Official != nil {
		update["official"] = *input.Official
	}

	err := p.queries.FilterUserID(userID).Update(ctx, update)
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
