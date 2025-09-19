package app

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"

	"github.com/chains-lab/profiles-svc/internal/app/domain/profiles"
	"github.com/chains-lab/profiles-svc/internal/config"
	"github.com/chains-lab/profiles-svc/internal/dbx"
	"github.com/chains-lab/profiles-svc/internal/errx"
)

type App struct {
	profiles profiles.Profiles

	db *sql.DB
}

func (a App) generateUsername() (string, error) {
	const (
		prefix = "citizen"
		digits = 8
	)
	buf := make([]byte, digits)
	if _, err := rand.Read(buf); err != nil {
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("cannot generate random digits: %w", err),
		)
	}
	for i := 0; i < digits; i++ {
		buf[i] = '0' + (buf[i] % 10)
	}
	return prefix + string(buf), nil
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	profiles, err := profiles.NewProfile(pg)
	if err != nil {
		return App{}, err
	}

	return App{
		profiles: profiles,

		db: pg,
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
