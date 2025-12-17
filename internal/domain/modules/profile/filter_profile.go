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
	Verified        *bool
}

func (s Service) FilterProfile(ctx context.Context, params FilterParams, offset, limit int32) (entity.ProfileCollection, error) {
	collection, err := s.db.FilterProfiles(ctx, params, uint(offset), uint(limit))
	if err != nil {
		return entity.ProfileCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting profile with username '%s': %w", *params.UsernamePrefix, err),
		)
	}

	return collection, nil
}
