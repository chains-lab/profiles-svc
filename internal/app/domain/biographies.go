package domain

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/elector-cab-svc/internal/dbx"
	"github.com/google/uuid"
)

type BiographiesQ interface {
	New() dbx.BiographiesQ

	Insert(ctx context.Context, input dbx.BioModel) error
	Update(ctx context.Context, input dbx.UpdateBioInput) error
	Select(ctx context.Context) ([]dbx.BioModel, error)
	Get(ctx context.Context) (dbx.BioModel, error)
	Delete(ctx context.Context) error

	FilterUserID(userID uuid.UUID) dbx.BiographiesQ

	Count(ctx context.Context) (int, error)
	Page(limit, offset uint64) dbx.BiographiesQ
}

type Biographies struct {
	queries BiographiesQ
}

func NewBiographies(db *sql.DB) (Biographies, error) {
	return Biographies{
		queries: dbx.NewBiographies(db),
	}, nil
}

func (b Biographies) Create(ctx context.Context, userID uuid.UUID) error {
	if err := b.queries.Insert(ctx, dbx.BioModel{
		UserID: userID,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func (b Biographies) Get(ctx context.Context, userID uuid.UUID) (models.Bio, error) {
	bio, err := b.queries.New().FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		default:
			return models.Bio{}, ape.ErrorInternal(err) //TODO
		}
	}

	return BioFromDb(bio), nil
}

func (b Biographies) SetSex(ctx context.Context, userID uuid.UUID, sex string) error {
	if !models.ValidateSex(sex) {
		return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
	}

	now := time.Now().UTC()

	bio, err := b.Get(ctx, userID)
	if err != nil {
		return err
	}

	if bio.SexUpdatedAt != nil {
		last := *bio.SexUpdatedAt

		if now.Sub(last) < 365*24*time.Hour {
			//nextAllowed := last.Add(365 * 24 * time.Hour)
			//wait := nextAllowed.Sub(now).Round(time.Hour)

			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}
	}

	if err := b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		Sex:          &sex,
		SexUpdatedAt: &now,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func (b Biographies) SetBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) error {
	bio, err := b.Get(ctx, userID)
	if err != nil {
		return err
	}

	if bio.Birthday != nil {
		return ape.ErrorInternal(fmt.Errorf("birthday is already set")) //TODO: add error
	}

	if err := b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		Birthday: &birthday,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func (b Biographies) SetCitizenship(ctx context.Context, userID uuid.UUID, citizenship string) error {
	bio, err := b.Get(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if bio.CitizenshipUpdatedAt != nil {
		last := *bio.CitizenshipUpdatedAt

		if now.Sub(last) < 365*24*time.Hour {
			//nextAllowed := last.Add(365 * 24 * time.Hour)
			//wait := nextAllowed.Sub(now).Round(time.Hour)

			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}
	}

	if err := b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		Citizenship:          &citizenship,
		CitizenshipUpdatedAt: &now,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func (b Biographies) SetNationality(ctx context.Context, userID uuid.UUID, nationality string) error {
	bio, err := b.Get(ctx, userID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()

	//TODO validate nationality from other api
	if bio.NationalityUpdatedAt != nil {
		last := *bio.NationalityUpdatedAt

		if now.Sub(last) < 365*24*time.Hour {
			//nextAllowed := last.Add(365 * 24 * time.Hour)
			//wait := nextAllowed.Sub(now).Round(time.Hour)

			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}
	}

	if err := b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		Nationality:          &nationality,
		NationalityUpdatedAt: &now,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func (b Biographies) SetPrimaryLanguage(ctx context.Context, userID uuid.UUID, primaryLanguage string) error {
	bio, err := b.Get(ctx, userID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()

	//TODO validate primaryLanguage from other api
	if bio.PrimaryLanguageUpdatedAt != nil {
		last := *bio.PrimaryLanguageUpdatedAt

		if now.Sub(last) < 365*24*time.Hour {
			//nextAllowed := last.Add(365 * 24 * time.Hour)
			//wait := nextAllowed.Sub(now).Round(time.Hour)

			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}
	}

	if err := b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		PrimaryLanguage:          &primaryLanguage,
		PrimaryLanguageUpdatedAt: &now,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

type AdminBioUpdate struct {
	Birthday        *time.Time
	Sex             *string
	Citizenship     *string
	Nationality     *string
	PrimaryLanguage *string
}

func (b Biographies) AdminUpdateBio(ctx context.Context, userID uuid.UUID, input AdminBioUpdate) error {
	_, err := b.Get(ctx, userID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()

	var dbInput dbx.UpdateBioInput

	if input.Birthday != nil {
		dbInput.Birthday = input.Birthday
	}
	if input.Sex != nil {
		if !models.ValidateSex(*input.Sex) {
			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}

		dbInput.Sex = input.Sex
		dbInput.SexUpdatedAt = &now
	}
	if input.Citizenship != nil {
		//TODO validate citizenship from other api
		dbInput.Citizenship = input.Citizenship
		dbInput.CitizenshipUpdatedAt = &now
	}
	if input.Nationality != nil {
		//TODO validate
		dbInput.Nationality = input.Nationality
		dbInput.NationalityUpdatedAt = &now
	}
	if input.PrimaryLanguage != nil {
		//TODO validate
		dbInput.PrimaryLanguage = input.PrimaryLanguage
		dbInput.PrimaryLanguageUpdatedAt = &now
	}

	if err := b.queries.New().FilterUserID(userID).Update(ctx, dbInput); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func BioFromDb(input dbx.BioModel) models.Bio {
	return models.Bio{
		UserID:          input.UserID,
		Sex:             input.Sex,
		Birthday:        input.Birthday,
		Citizenship:     input.Citizenship,
		Nationality:     input.Nationality,
		PrimaryLanguage: input.PrimaryLanguage,

		SexUpdatedAt:             input.SexUpdatedAt,
		NationalityUpdatedAt:     input.NationalityUpdatedAt,
		CitizenshipUpdatedAt:     input.CitizenshipUpdatedAt,
		PrimaryLanguageUpdatedAt: input.PrimaryLanguageUpdatedAt,
	}
}
