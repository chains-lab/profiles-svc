package app

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/app/domain"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) ProfileGetUserID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	return a.profiles.GetByID(ctx, userID)
}

func (a App) ProfileGetUsername(ctx context.Context, username string) (models.Profile, error) {
	return a.profiles.GetByUsername(ctx, username)
}

type UpdateProfileInput struct {
	Username    *string `json:"username,omitempty"`
	Pseudonym   *string `json:"pseudonym,omitempty"`
	Description *string `json:"description,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
	Official    *bool   `json:"official,omitempty"`
}

func (a App) ProfileUpdate(ctx context.Context, userID uuid.UUID, profile UpdateProfileInput) (models.Profile, error) {
	err := a.profiles.Update(ctx, userID, entities.UpdateProfileInput{
		Username:    profile.Username,
		Pseudonym:   profile.Pseudonym,
		Description: profile.Description,
		Avatar:      profile.Avatar,
		Official:    profile.Official,
	})
	if err != nil {
		return models.Profile{}, err
	}

	return a.ProfileGetUserID(ctx, userID)
}
