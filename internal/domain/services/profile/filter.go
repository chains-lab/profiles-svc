package profile

import (
	"context"
	"fmt"

	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/chains-lab/profiles-svc/internal/domain/models"
	"github.com/google/uuid"
)

type FilterParams struct {
	UserID       []uuid.UUID
	Username     []string
	UsernameLike *string
	Pseudonym    *string
	Official     *bool
}

func (s Service) Filter(ctx context.Context, filters FilterParams, page uint64, size uint64) (models.ProfileCollection, error) {
	res, err := s.db.FilterProfiles(ctx, filters, page, size)
	if err != nil {
		return models.ProfileCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter profile: %w", err),
		)
	}

	return res, nil
}
