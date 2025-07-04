package app

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/chains-lab/elector-cab-svc/internal/app/entities"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	return a.profiles.GetByID(ctx, userID)
}

func (a App) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	return a.profiles.GetByUsername(ctx, username)
}

type UpdateProfileInput struct {
	Pseudonym   *string `json:"pseudonym,omitempty"`
	Description *string `json:"description,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
}

func (a App) UpdateProfile(ctx context.Context, userID uuid.UUID, profile UpdateProfileInput) (models.Profile, error) {
	err := a.profiles.Update(ctx, userID, entities.UpdateProfileInput{
		Pseudonym:   profile.Pseudonym,
		Description: profile.Description,
		Avatar:      profile.Avatar,
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

type AdminUpdateProfileInput struct {
	Pseudonym   *string `json:"pseudonym,omitempty"`
	Description *string `json:"description,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
	Official    *bool   `json:"official,omitempty"`
}

func (a App) AdminUpdateProfile(ctx context.Context, userID uuid.UUID, profile AdminUpdateProfileInput) (models.Profile, error) {
	prof, err := a.GetProfileByUserID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	err = a.profiles.Update(ctx, userID, entities.UpdateProfileInput{
		Pseudonym:   profile.Pseudonym,
		Description: profile.Description,
		Avatar:      profile.Avatar,
		Official:    profile.Official,
	})
	if err != nil {
		return models.Profile{}, err
	}

	return prof, nil
}

func (a App) ResetUsername(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	prof, err := a.GetProfileByUserID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	generateUsername := func() (string, error) {
		const (
			prefix = "elector"
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

	return prof, nil
}

type ResetUserProfileInput struct {
	Pseudonym   bool `json:"pseudonym"`
	Description bool `json:"description"`
	Avatar      bool `json:"avatar"`
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
