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

type ResidencesQ interface {
	New() dbx.ResidencesQ

	Insert(ctx context.Context, input dbx.ResidenceModel) error
	Update(ctx context.Context, input dbx.UpdateResidenceInput) error
	Get(ctx context.Context) (dbx.ResidenceModel, error)
	Select(ctx context.Context) ([]dbx.ResidenceModel, error)
	Delete(ctx context.Context) error

	FilterUserID(userID uuid.UUID) dbx.ResidencesQ

	Page(limit, offset uint64) dbx.ResidencesQ
	Count(ctx context.Context) (int, error)
}

type Residences struct {
	queries ResidencesQ
}

func NewResidences(db *sql.DB) (Residences, error) {
	return Residences{
		queries: dbx.NewResidences(db),
	}, nil
}

func (r Residences) Create(ctx context.Context, userID uuid.UUID) error {
	if err := r.queries.Insert(ctx, dbx.ResidenceModel{
		UserID: userID,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(fmt.Errorf("")) //TODO: handle error properly
		}
	}

	return nil
}

func (r Residences) Get(ctx context.Context, userID uuid.UUID) (models.Residence, error) {
	residence, err := r.queries.FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		default:
			return models.Residence{}, ape.ErrorInternal(err) //TODO: handle error properly
		}
	}

	return ResidenceFromDb(residence), nil
}

func (r Residences) Update(ctx context.Context, userID uuid.UUID, country string, city string) error {
	residence, err := r.Get(ctx, userID)
	if err != nil {
		return ape.ErrorInternal(err) //TODO: handle error properly
	}

	now := time.Now().UTC()

	//TODO  Validate residence
	if residence.UpdatedAt != nil {
		last := *residence.UpdatedAt

		if now.Sub(last) < 365*24*time.Hour {
			return ape.ErrorInternal(fmt.Errorf(""))
		}
	}

	if err := r.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateResidenceInput{
		City:      &city,
		Country:   &country,
		UpdatedAt: &now,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO: handle error properly
		}
	}

	return nil
}

func (r Residences) UpdateAdmin(ctx context.Context, userID uuid.UUID, country string, city string) error {
	_, err := r.Get(ctx, userID)
	if err != nil {
		return ape.ErrorInternal(err) //TODO: handle error properly
	}

	now := time.Now().UTC()

	//TODO  Validate residence

	if err := r.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateResidenceInput{
		City:      &city,
		Country:   &country,
		UpdatedAt: &now,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO: handle error properly
		}
	}

	return nil
}

func ResidenceFromDb(model dbx.ResidenceModel) models.Residence {
	return models.Residence{
		UserID:    model.UserID,
		Country:   model.Country,
		City:      model.City,
		UpdatedAt: model.UpdatedAt,
	}
}
