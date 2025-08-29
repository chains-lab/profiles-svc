package app

import (
	"context"
	"time"

	"github.com/chains-lab/profiles-svc/internal/app/entity"
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

func (a App) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	return a.profiles.GetByID(ctx, userID)
}

func (a App) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	return a.profiles.GetByUsername(ctx, username)
}

type UpdateProfileInput struct {
	Pseudonym   *string
	Description *string
	Avatar      *string
	Sex         *string
	BirthDate   *time.Time
}

func (a App) UpdateProfile(ctx context.Context, userID uuid.UUID, profile UpdateProfileInput) (models.Profile, error) {
	err := a.profiles.Update(ctx, userID, entity.UpdateProfileInput{
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

	if err = a.profiles.Update(ctx, userID, entity.UpdateProfileInput{
		Official: &official,
	}); err != nil {
		return models.Profile{}, err
	}

	return a.GetProfileByUserID(ctx, userID)
}

func (a App) ResetUsername(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	_, err := a.GetProfileByUserID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	usrnm, err := a.profiles.GenerateUsername()
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
	dmInput := entity.UpdateProfileInput{}
	dmInput.Pseudonym = &empty
	dmInput.Description = &empty
	dmInput.Avatar = &empty

	if err = a.profiles.Update(ctx, userID, dmInput); err != nil {
		return models.Profile{}, err
	}

	return a.GetProfileByUserID(ctx, userID)
}
