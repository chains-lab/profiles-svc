package app

import (
	"context"
	"time"

	"github.com/chains-lab/profiles-svc/internal/app/domain/profiles"
	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/google/uuid"
)

type UpdateProfileInput struct {
	Pseudonym   *string
	Description *string
	Avatar      *string
	Sex         *string
	BirthDate   *time.Time
}

func (a App) UpdateProfile(ctx context.Context, userID uuid.UUID, profile UpdateProfileInput) (models.Profile, error) {
	err := a.profiles.Update(ctx, userID, profiles.UpdateProfileInput{
		Pseudonym:   profile.Pseudonym,
		Description: profile.Description,
		Avatar:      profile.Avatar,
		Sex:         profile.Sex,
		BirthDate:   profile.BirthDate,
	})
	if err != nil {
		return models.Profile{}, err
	}

	return a.GetProfileByUserID(ctx, userID)
}

func (a App) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error) {
	if err := a.profiles.UpdateUsername(ctx, userID, username); err != nil {
		return models.Profile{}, err
	}

	return a.GetProfileByUserID(ctx, userID)
}

func (a App) UpdateProfileOfficial(ctx context.Context, userID uuid.UUID, official bool) (models.Profile, error) {
	prof, err := a.GetProfileByUserID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	if prof.Official == official {
		return prof, nil
	}

	if err = a.profiles.Update(ctx, userID, profiles.UpdateProfileInput{
		Official: &official,
	}); err != nil {
		return models.Profile{}, err
	}

	return a.GetProfileByUserID(ctx, userID)
}
