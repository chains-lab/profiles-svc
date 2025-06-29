package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/elector-cab-svc/internal/app/domain"
	"github.com/chains-lab/elector-cab-svc/internal/dbx"
	"github.com/chains-lab/elector-cab-svc/internal/utils/config"
)

type App struct {
	profiles    entities.Profiles
	jobResumes  entities.JobResumes
	biographies entities.Biographies

	db *sql.DB
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	profiles, err := entities.NewProfile(pg)
	if err != nil {
		return App{}, err
	}
	jobs, err := entities.NewJobResumes(pg)
	if err != nil {
		return App{}, err
	}
	biographies, err := entities.NewBiographies(pg)
	if err != nil {
		return App{}, err
	}

	return App{
		profiles:    profiles,
		jobResumes:  jobs,
		biographies: biographies,
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
