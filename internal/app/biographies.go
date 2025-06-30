package app

import (
	"context"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/app/entities"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) UpdateResidence(ctx context.Context, userID uuid.UUID, city, country string) (models.Biography, error) {
	err := a.biographies.UpdateResidence(ctx, userID, city, country)
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}

func (a App) AdminUpdateResidence(ctx context.Context, userID uuid.UUID, city, country string) (models.Biography, error) {
	err := a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		Country: &country,
		City:    &city,
	})
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}

func (a App) UpdateSex(ctx context.Context, userID uuid.UUID, sex string) (models.Biography, error) {
	err := a.biographies.UpdateSex(ctx, userID, sex)
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}

func (a App) AdminUpdateSex(ctx context.Context, userID uuid.UUID, sex string) (models.Biography, error) {
	err := a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		Sex: &sex,
	})
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}

func (a App) UpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) (models.Biography, error) {
	err := a.biographies.UpdateBirthday(ctx, userID, birthday)
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}

func (a App) AdminUpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) (models.Biography, error) {
	err := a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		Birthday: &birthday,
	})
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}

func (a App) UpdateNationality(ctx context.Context, userID uuid.UUID, nationality string) (models.Biography, error) {
	err := a.biographies.SetNationality(ctx, userID, nationality)
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}

func (a App) AdminUpdateNationality(ctx context.Context, userID uuid.UUID, nationality string) (models.Biography, error) {
	err := a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		Nationality: &nationality,
	})
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}

func (a App) UpdatePrimaryLanguage(ctx context.Context, userID uuid.UUID, primaryLanguage string) (models.Biography, error) {
	err := a.biographies.SetPrimaryLanguage(ctx, userID, primaryLanguage)
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}

func (a App) AdminUpdatePrimaryLanguage(ctx context.Context, userID uuid.UUID, primaryLanguage string) (models.Biography, error) {
	err := a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		PrimaryLanguage: &primaryLanguage,
	})
	if err != nil {
		return models.Biography{}, err
	}

	return a.biographies.GetByUserID(ctx, userID)
}
