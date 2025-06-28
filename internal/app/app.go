package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/elector-cab-svc/internal/app/domain"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/elector-cab-svc/internal/dbx"
	"github.com/chains-lab/elector-cab-svc/internal/utils/config"
	"github.com/google/uuid"
)

type App struct {
	profiles    domain.Profiles
	jobs        domain.Jobs
	biographies domain.Biographies
	residences  domain.Residences

	db *sql.DB
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	profiles, err := domain.NewProfile(pg)
	if err != nil {
		return App{}, err
	}
	jobs, err := domain.NewJob(pg)
	if err != nil {
		return App{}, err
	}
	biographies, err := domain.NewBiographies(pg)
	if err != nil {
		return App{}, err
	}
	residences, err := domain.NewResidences(pg)
	if err != nil {
		return App{}, err
	}

	return App{
		profiles:    profiles,
		jobs:        jobs,
		biographies: biographies,
		residences:  residences,
		db:          pg,
	}, nil
}

func (a App) transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, dbx.TxKey, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a App) UserDataGetUserID(ctx context.Context, userID uuid.UUID) (models.UserData, error) {
	profile, err := a.profiles.GetByID(ctx, userID)
	if err != nil {
		return models.UserData{}, err
	}

	jobs, err := a.jobs.Get(ctx, userID)
	if err != nil {
		return models.UserData{}, err
	}

	biography, err := a.biographies.Get(ctx, userID)
	if err != nil {
		return models.UserData{}, err
	}

	residence, err := a.residences.Get(ctx, userID)
	if err != nil {
		return models.UserData{}, err
	}

	return models.UserData{
		Profile:   profile,
		Job:       jobs,
		Bio:       biography,
		Residence: residence,
	}, nil
}

func (a App) UserDataGetByUsername(ctx context.Context, username string) (models.UserData, error) {
	profile, err := a.profiles.GetByUsername(ctx, username)
	if err != nil {
		return models.UserData{}, err
	}

	jobs, err := a.jobs.Get(ctx, profile.UserID)
	if err != nil {
		return models.UserData{}, err
	}

	biography, err := a.biographies.Get(ctx, profile.UserID)
	if err != nil {
		return models.UserData{}, err
	}

	residence, err := a.residences.Get(ctx, profile.UserID)
	if err != nil {
		return models.UserData{}, err
	}

	return models.UserData{
		Profile:   profile,
		Job:       jobs,
		Bio:       biography,
		Residence: residence,
	}, nil
}

func (a App) GetUserBiography(ctx context.Context, userID uuid.UUID) (models.Bio, error) {
	return a.biographies.Get(ctx, userID)
}

func (a App) GetUserResidence(ctx context.Context, userID uuid.UUID) (models.Residence, error) {
	return a.residences.Get(ctx, userID)
}

func (a App) GetUserJob(ctx context.Context, userID uuid.UUID) (models.Job, error) {
	return a.jobs.Get(ctx, userID)
}
