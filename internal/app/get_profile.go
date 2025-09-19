package app

import (
	"context"

	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	return a.profiles.GetByID(ctx, userID)
}

func (a App) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	return a.profiles.GetByUsername(ctx, username)
}
