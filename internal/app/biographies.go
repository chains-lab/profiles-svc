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

func (a App) UpdateResidence(ctx context.Context, userID uuid.UUID, input UpdateResidenceInput) (models.Biography, error) {
	res, err := a.biographies.UpdateResidence(ctx, userID, entities.UpdateResidenceInput{
		City:    input.City,
		Region:  input.Region,
		Country: input.Country,
	})
	if err != nil {
		return models.Biography{}, err
	}

	return res, nil
}

func (a App) UpdateSex(ctx context.Context, userID uuid.UUID, sex string) (models.Biography, error) {
	res, err := a.biographies.UpdateSex(ctx, userID, sex)
	if err != nil {
		return models.Biography{}, err
	}

	return res, nil
}

func (a App) UpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) (models.Biography, error) {
	res, err := a.biographies.UpdateBirthday(ctx, userID, birthday)
	if err != nil {
		return models.Biography{}, err
	}

	return res, nil
}

type UpdateBiographyInput struct {
	Birthday *time.Time
	Sex      *string
	City     *string
	Region   *string
	Country  *string
}

func (a App) AdminUpdateBiography(ctx context.Context, userID uuid.UUID, input UpdateBiographyInput) (models.Biography, error) {
	res, err := a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		Birthday: input.Birthday,
		Sex:      input.Sex,
		City:     input.City,
		Region:   input.Region,
		Country:  input.Country,
	})
	if err != nil {
		return models.Biography{}, err
	}

	return res, nil
}
