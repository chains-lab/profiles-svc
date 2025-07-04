package app

import (
	"context"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/app/entities"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/google/uuid"
)

type UpdateResidenceInput struct {
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
}

func (a App) UpdateResidence(ctx context.Context, userID uuid.UUID, input UpdateResidenceInput) error {
	err := a.biographies.UpdateResidence(ctx, userID, entities.UpdateResidenceInput{
		City:    input.City,
		Region:  input.Region,
		Country: input.Country,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a App) UpdateSex(ctx context.Context, userID uuid.UUID, sex string) error {
	err := a.biographies.UpdateSex(ctx, userID, sex)
	if err != nil {
		return err
	}

	return nil
}

func (a App) UpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) error {
	err := a.biographies.UpdateBirthday(ctx, userID, birthday)
	if err != nil {
		return err
	}

	return nil
}

func (a App) UpdateNationality(ctx context.Context, userID uuid.UUID, nationality string) error {
	err := a.biographies.SetNationality(ctx, userID, nationality)
	if err != nil {
		return err
	}

	return nil
}

func (a App) UpdatePrimaryLanguage(ctx context.Context, userID uuid.UUID, primaryLanguage string) error {
	err := a.biographies.SetPrimaryLanguage(ctx, userID, primaryLanguage)
	if err != nil {
		return err
	}

	return nil
}

type UpdateBiographyInput struct {
	Birthday        *time.Time
	Sex             *string
	City            *string
	Region          *string
	Country         *string
	Nationality     *string
	PrimaryLanguage *string
}

func (a App) AdminUpdateBiography(ctx context.Context, userID uuid.UUID, input UpdateBiographyInput) (models.Biography, error) {
	err := a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		Birthday:        input.Birthday,
		Sex:             input.Sex,
		City:            input.City,
		Country:         input.Country,
		Nationality:     input.Nationality,
		PrimaryLanguage: input.PrimaryLanguage,
	})
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}
