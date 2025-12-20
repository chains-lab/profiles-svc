package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/profiles-svc/internal/domain/entity"
	"github.com/umisto/profiles-svc/internal/domain/errx"
)

func (s Service) CreateProfile(ctx context.Context, userID uuid.UUID, username string) (entity.Profile, error) {
	profile, err := s.db.CreateProfile(ctx, userID, username)
	if err != nil {
		return entity.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("creating profile for user '%s': %w", userID, err),
		)
	}

	return profile, nil
}
