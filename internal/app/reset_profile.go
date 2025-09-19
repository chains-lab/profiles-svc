package app

import (
	"context"

	"github.com/chains-lab/profiles-svc/internal/app/domain/profiles"
	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) ResetUsername(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	_, err := a.GetProfileByUserID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	usrnm, err := a.generateUsername()
	if err != nil {
		return models.Profile{}, err
	}

	if err := a.profiles.UpdateUsername(ctx, userID, usrnm); err != nil {
		return models.Profile{}, err
	}

	return a.GetProfileByUserID(ctx, userID)
}

func (a App) ResetUserProfile(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	_, err := a.GetProfileByUserID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	empty := ""
	dmInput := profiles.UpdateProfileInput{}
	dmInput.Pseudonym = &empty
	dmInput.Description = &empty
	dmInput.Avatar = &empty

	if err = a.profiles.Update(ctx, userID, dmInput); err != nil {
		return models.Profile{}, err
	}

	return a.GetProfileByUserID(ctx, userID)
}
