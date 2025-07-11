package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/ape"
	"github.com/chains-lab/elector-cab-svc/internal/app/domain"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/elector-cab-svc/internal/app/references"
	"github.com/chains-lab/elector-cab-svc/internal/config"
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
	cfg     config.Config
}

func NewBiographies(db *sql.DB, cfg config.Config) (Biographies, error) {
	return Biographies{
		queries: dbx.NewBiographies(db),
		cfg:     cfg,
	}, nil
}

func (b Biographies) Create(ctx context.Context, userID uuid.UUID) error {
	if err := b.queries.Insert(ctx, dbx.BioModel{
		UserID: userID,
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseProfileForUserAlreadyExists(err, userID)
		default:
			return ape.RaiseInternal(err)
		}
	}

	return nil
}

func (b Biographies) GetByUserID(ctx context.Context, userID uuid.UUID) (models.Biography, error) {
	bio, err := b.queries.New().FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Biography{}, ape.RaiseProfileForUserNotFound(err, userID)
		default:
			return models.Biography{}, ape.RaiseInternal(err)
		}
	}

	return BioFromDb(bio), nil
}

func (b Biographies) UpdateSex(ctx context.Context, userID uuid.UUID, sex string) (models.Biography, error) {
	if err := references.ValidateSex(sex); err != nil {
		return models.Biography{}, ape.RaiseSexIsNotValid(err)
	}

	now := time.Now().UTC()

	bio, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return models.Biography{}, err
	}

	if bio.SexUpdatedAt != nil {
		last := *bio.SexUpdatedAt

		if err := domain.ValidateUpdateProperty(last, 365*24*time.Hour); err != nil {
			return models.Biography{}, ape.RaiseSexUpdateCooldown(err, userID)
		}
	}

	if err = b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		Sex:          &sex,
		SexUpdatedAt: &now,
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Biography{}, ape.RaiseProfileForUserNotFound(err, userID)
		default:
			return models.Biography{}, ape.RaiseInternal(err)
		}
	}

	return models.Biography{
		UserID:             userID,
		Birthday:           bio.Birthday,
		Sex:                &sex,
		City:               bio.City,
		Region:             bio.Region,
		Country:            bio.Country,
		SexUpdatedAt:       &now,
		ResidenceUpdatedAt: bio.ResidenceUpdatedAt,
	}, nil
}

func (b Biographies) UpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) (models.Biography, error) {
	bio, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return models.Biography{}, err
	}

	if bio.Birthday != nil {
		return models.Biography{}, ape.RaiseBirthdayIsAlreadySet(fmt.Errorf("birthday is already set"), userID)
	}

	if err = b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		Birthday: &birthday,
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Biography{}, ape.RaiseProfileForUserNotFound(err, userID)
		default:
			return models.Biography{}, ape.RaiseInternal(err) //TODO
		}
	}

	return models.Biography{
		UserID:             userID,
		Birthday:           &birthday,
		Sex:                bio.Sex,
		City:               bio.City,
		Region:             bio.Region,
		Country:            bio.Country,
		SexUpdatedAt:       bio.SexUpdatedAt,
		ResidenceUpdatedAt: bio.ResidenceUpdatedAt,
	}, nil
}

type UpdateResidenceInput struct {
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
}

func (b Biographies) UpdateResidence(ctx context.Context, userID uuid.UUID, req UpdateResidenceInput) (models.Biography, error) {
	//TODO validate country and city from other api
	err := references.ValidateResidence(req.City, req.Region, req.Country, b.cfg.Properties.Residence.ApiKey)
	if err != nil {
		return models.Biography{}, ape.RaiseResidenceIsNotValid(err)
	}

	bio, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return models.Biography{}, err
	}

	now := time.Now().UTC()

	if bio.ResidenceUpdatedAt != nil {
		last := *bio.ResidenceUpdatedAt

		if err := domain.ValidateUpdateProperty(last, 365*24*time.Hour); err != nil {
			return models.Biography{}, ape.RaiseResidenceUpdateCooldown(err, userID)
		}
	}

	if err = b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		City:               &req.City,
		Region:             &req.Region,
		Country:            &req.Country,
		ResidenceUpdatedAt: &now,
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Biography{}, ape.RaiseProfileForUserNotFound(err, userID)
		default:
			return models.Biography{}, ape.RaiseInternal(err)
		}
	}

	return models.Biography{
		UserID:             userID,
		Birthday:           bio.Birthday,
		Sex:                bio.Sex,
		City:               &req.City,
		Region:             &req.Region,
		Country:            &req.Country,
		SexUpdatedAt:       bio.SexUpdatedAt,
		ResidenceUpdatedAt: &now,
	}, nil
}

type AdminBioUpdate struct {
	Birthday *time.Time
	Sex      *string
	City     *string
	Region   *string
	Country  *string
}

func (b Biographies) AdminUpdateBio(ctx context.Context, userID uuid.UUID, input AdminBioUpdate) (models.Biography, error) {
	u, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return models.Biography{}, err
	}

	now := time.Now().UTC()

	var dbInput dbx.UpdateBioInput

	if input.Birthday != nil {
		if u.Birthday != nil {
			return models.Biography{}, ape.RaiseBirthdayIsAlreadySet(fmt.Errorf("birthday is already set"), userID)
		}

		if input.Birthday.Year() < 1900 {
			return models.Biography{}, ape.RaiseBirthdayIsNotValid(fmt.Errorf("birthday is not valid, must be after 1900"))
		}

		dbInput.Birthday = input.Birthday
	}

	if input.Sex != nil {
		if err := references.ValidateSex(*input.Sex); err != nil {
			return models.Biography{}, ape.RaiseSexIsNotValid(err)
		}

		dbInput.Sex = input.Sex
		dbInput.SexUpdatedAt = &now
	}

	if input.City != nil && input.Country != nil {

		//TODO check functionality of this validation
		if err = references.ValidateResidence(*input.City, *input.Region, *input.Country, b.cfg.Properties.Residence.ApiKey); err != nil {
			return models.Biography{}, ape.RaiseResidenceIsNotValid(err)
		}

		dbInput.City = input.City
		dbInput.Country = input.Country
		dbInput.Region = input.Region
		dbInput.ResidenceUpdatedAt = &now
	}

	if err = b.queries.New().FilterUserID(userID).Update(ctx, dbInput); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Biography{}, ape.RaiseProfileForUserNotFound(err, userID)
		default:
			return models.Biography{}, ape.RaiseInternal(err)
		}
	}

	bio, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return models.Biography{}, err
	}

	return bio, nil
}

func BioFromDb(input dbx.BioModel) models.Biography {
	return models.Biography{
		UserID:   input.UserID,
		Birthday: input.Birthday,
		Sex:      input.Sex,
		City:     input.City,
		Region:   input.Region,
		Country:  input.Country,

		SexUpdatedAt:       input.SexUpdatedAt,
		ResidenceUpdatedAt: input.ResidenceUpdatedAt,
	}
}
