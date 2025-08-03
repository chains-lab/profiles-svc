package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/chains-lab/profiles-svc/internal/app/entities"
	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/google/uuid"
)

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
	err := a.profiles.Update(ctx, userID, entities.UpdateProfileInput{
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

func (a App) AdminUpdateProfileOfficial(ctx context.Context, userID uuid.UUID, official bool) (models.Profile, error) {
	prof, err := a.GetProfileByUserID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	if prof.Official == official {
		return prof, nil
	}

	if err = a.profiles.Update(ctx, userID, entities.UpdateProfileInput{
		Official: &official,
	}); err != nil {
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

func (a App) ResetUsername(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	_, err := a.GetProfileByUserID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	generateUsername := func() (string, error) {
		const (
			prefix = "citizen"
			digits = 8
		)
		buf := make([]byte, digits)
		if _, err := rand.Read(buf); err != nil {
			return "", fmt.Errorf("cannot generate random digits: %w", err)
		}
		for i := 0; i < digits; i++ {
			buf[i] = '0' + (buf[i] % 10)
		}
		return prefix + string(buf), nil
	}

	username, err := generateUsername()
	if err != nil {
		return models.Profile{}, err
	}

	if err := a.profiles.UpdateUsername(ctx, userID, username); err != nil {
		return models.Profile{}, err
	}

	return a.GetProfileByUserID(ctx, userID)
}

type ResetUserProfileInput struct {
	Pseudonym   bool
	Description bool
	Avatar      bool
}

func (a App) ResetUserProfile(ctx context.Context, userID uuid.UUID, input ResetUserProfileInput) (models.Profile, error) {
	_, err := a.GetProfileByUserID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	empty := ""

	dmInput := entities.UpdateProfileInput{}
	if input.Pseudonym {
		dmInput.Pseudonym = &empty
	}
	if input.Description {
		dmInput.Description = &empty
	}
	if input.Avatar {
		dmInput.Avatar = &empty
	}

	if err := a.profiles.Update(ctx, userID, dmInput); err != nil {
		return models.Profile{}, err
	}

	res, err := a.GetProfileByUserID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	return res, nil
}

type CreateProfileInput struct {
	Username    string
	Pseudonym   *string
	Description *string
	Avatar      *string
	Sex         *string
	BirthDate   *time.Time
}

func (a App) CreateProfile(ctx context.Context, userID uuid.UUID, input CreateProfileInput) (models.Profile, error) {
	err := a.profiles.Create(ctx, userID, entities.CreateProfileInput{
		Username:    input.Username,
		Pseudonym:   input.Pseudonym,
		Description: input.Description,
		Avatar:      input.Avatar,
		Sex:         input.Sex,
		BirthDate:   input.BirthDate,
	})
	if err != nil {
		return models.Profile{}, err
	}

	return a.GetProfileByUserID(ctx, userID)
}
