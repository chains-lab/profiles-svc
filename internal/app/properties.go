package app

import (
	"context"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/app/domain"
	"github.com/google/uuid"
)

func (a App) UpdateResidence(ctx context.Context, userID uuid.UUID, city, country string) error {
	return a.residences.Update(ctx, userID, city, country)
}

func (a App) UpdateResidenceAdmin(ctx context.Context, userID uuid.UUID, city, country string) error {
	return a.residences.UpdateAdmin(ctx, userID, city, country)
}

func (a App) UpdateSex(ctx context.Context, userID uuid.UUID, sex string) error {
	return a.biographies.SetSex(ctx, userID, sex)
}

func (a App) AdminUpdateSex(ctx context.Context, userID uuid.UUID, sex string) error {
	return a.biographies.AdminUpdateBio(ctx, userID, domain.AdminBioUpdate{
		Sex: &sex,
	})
}

func (a App) UpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) error {
	return a.biographies.SetBirthday(ctx, userID, birthday)
}

func (a App) AdminUpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) error {
	return a.biographies.AdminUpdateBio(ctx, userID, domain.AdminBioUpdate{
		Birthday: &birthday,
	})
}

func (a App) UpdateCitizenship(ctx context.Context, userID uuid.UUID, citizenship string) error {
	return a.biographies.SetCitizenship(ctx, userID, citizenship)
}

func (a App) AdminUpdateCitizenship(ctx context.Context, userID uuid.UUID, citizenship string) error {
	return a.biographies.AdminUpdateBio(ctx, userID, domain.AdminBioUpdate{
		Citizenship: &citizenship,
	})
}

func (a App) UpdateNationality(ctx context.Context, userID uuid.UUID, nationality string) error {
	return a.biographies.SetNationality(ctx, userID, nationality)
}

func (a App) AdminUpdateNationality(ctx context.Context, userID uuid.UUID, nationality string) error {
	return a.biographies.AdminUpdateBio(ctx, userID, domain.AdminBioUpdate{
		Nationality: &nationality,
	})
}

func (a App) UpdatePrimaryLanguage(ctx context.Context, userID uuid.UUID, primaryLanguage string) error {
	return a.biographies.SetPrimaryLanguage(ctx, userID, primaryLanguage)
}

func (a App) AdminUpdatePrimaryLanguage(ctx context.Context, userID uuid.UUID, primaryLanguage string) error {
	return a.biographies.AdminUpdateBio(ctx, userID, domain.AdminBioUpdate{
		PrimaryLanguage: &primaryLanguage,
	})
}

func (a App) UpdateDegree(ctx context.Context, userID uuid.UUID, degree string) error {
	return a.jobs.UpdateDegree(ctx, userID, degree)
}

func (a App) AdminUpdateDegree(ctx context.Context, userID uuid.UUID, degree string) error {
	return a.jobs.AdminUpdate(ctx, userID, domain.AdminJobUpdate{
		Degree: &degree,
	})
}

func (a App) UpdateIndustry(ctx context.Context, userID uuid.UUID, industry string) error {
	return a.jobs.UpdateIndustry(ctx, userID, industry)
}

func (a App) AdminUpdateIndustry(ctx context.Context, userID uuid.UUID, industry string) error {
	return a.jobs.AdminUpdate(ctx, userID, domain.AdminJobUpdate{
		Industry: &industry,
	})
}

func (a App) UpdateIncome(ctx context.Context, userID uuid.UUID, income string) error {
	return a.jobs.UpdateIncome(ctx, userID, income)
}

func (a App) AdminUpdateIncome(ctx context.Context, userID uuid.UUID, income string) error {
	return a.jobs.AdminUpdate(ctx, userID, domain.AdminJobUpdate{
		Income: &income,
	})
}
