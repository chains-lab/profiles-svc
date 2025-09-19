package app

import (
	"context"

	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) CreateProfile(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error) {
	err := a.profiles.Create(ctx, userID, username)
	if err != nil {
		return models.Profile{}, err
	}

	return a.GetProfileByUserID(ctx, userID)
}
