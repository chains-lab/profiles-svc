package app

import (
	"context"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/app/domain"
	"github.com/google/uuid"
)

func (a App) UpdateResidence(ctx context.Context, userID uuid.UUID, city, country string) error {
	return a.biographies.UpdateResidence(ctx, userID, city, country)
}

func (a App) UpdateResidenceAdmin(ctx context.Context, userID uuid.UUID, city, country string) error {
	return a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		Country: &country,
		City:    &city,
	})
}

func (a App) UpdateSex(ctx context.Context, userID uuid.UUID, sex string) error {
	return a.biographies.UpdateSex(ctx, userID, sex)
}

func (a App) AdminUpdateSex(ctx context.Context, userID uuid.UUID, sex string) error {
	return a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		Sex: &sex,
	})
}

func (a App) UpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) error {
	return a.biographies.UpdateBirthday(ctx, userID, birthday)
}

func (a App) AdminUpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) error {
	return a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		Birthday: &birthday,
	})
}

func (a App) UpdateNationality(ctx context.Context, userID uuid.UUID, nationality string) error {
	return a.biographies.SetNationality(ctx, userID, nationality)
}

func (a App) AdminUpdateNationality(ctx context.Context, userID uuid.UUID, nationality string) error {
	return a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		Nationality: &nationality,
	})
}

func (a App) UpdatePrimaryLanguage(ctx context.Context, userID uuid.UUID, primaryLanguage string) error {
	return a.biographies.SetPrimaryLanguage(ctx, userID, primaryLanguage)
}

func (a App) AdminUpdatePrimaryLanguage(ctx context.Context, userID uuid.UUID, primaryLanguage string) error {
	return a.biographies.AdminUpdateBio(ctx, userID, entities.AdminBioUpdate{
		PrimaryLanguage: &primaryLanguage,
	})
}
