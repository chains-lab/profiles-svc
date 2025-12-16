package profile

import (
	"context"
	"fmt"

	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/chains-lab/profiles-svc/internal/domain/errx"
)

type FilterParams struct {
	UsernamePrefix  *string
	PseudonymPrefix *string
}

func (s Service) FilterProfile(ctx context.Context, params FilterParams, offset, limit int32) (entity.ProfileCollection, error) {
	var collection entity.ProfileCollection
	var err error

	switch {
	case params.UsernamePrefix != nil:
		collection, err = s.db.FilterProfilesByUsername(ctx, *params.UsernamePrefix, offset, limit)
		if err != nil {
			return entity.ProfileCollection{}, errx.ErrorInternal.Raise(
				fmt.Errorf("getting profile with username '%s': %w", *params.UsernamePrefix, err),
			)
		}
	case params.PseudonymPrefix != nil:
		collection, err = s.db.FilterProfilesByPseudonym(ctx, *params.PseudonymPrefix, offset, limit)
		if err != nil {
			return entity.ProfileCollection{}, errx.ErrorInternal.Raise(
				fmt.Errorf("getting profile with pseudonym '%s': %w", *params.PseudonymPrefix, err),
			)
		}
	case params.UsernamePrefix == nil && params.PseudonymPrefix == nil:
		return entity.ProfileCollection{}, nil
	}

	return collection, nil
}
