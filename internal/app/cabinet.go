package app

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) CabinetGetByUserID(ctx context.Context, userID uuid.UUID) (models.Cabinet, error) {
	profile, err := a.profiles.GetByID(ctx, userID)
	if err != nil {
		return models.Cabinet{}, err
	}

	jobs, err := a.jobs.Get(ctx, userID)
	if err != nil {
		return models.Cabinet{}, err
	}

	biography, err := a.biographies.Get(ctx, userID)
	if err != nil {
		return models.Cabinet{}, err
	}

	return models.Cabinet{
		Profile:   profile,
		Job:       jobs,
		Biography: biography,
	}, nil
}

func (a App) CabinetGetByUsername(ctx context.Context, username string) (models.Cabinet, error) {
	profile, err := a.profiles.GetByUsername(ctx, username)
	if err != nil {
		return models.Cabinet{}, err
	}

	jobs, err := a.jobs.Get(ctx, profile.UserID)
	if err != nil {
		return models.Cabinet{}, err
	}

	biography, err := a.biographies.Get(ctx, profile.UserID)
	if err != nil {
		return models.Cabinet{}, err
	}

	return models.Cabinet{
		Profile:   profile,
		Job:       jobs,
		Biography: biography,
	}, nil
}

func (a App) GetUserBiography(ctx context.Context, userID uuid.UUID) (models.Biography, error) {
	return a.biographies.Get(ctx, userID)
}

func (a App) GetUserJob(ctx context.Context, userID uuid.UUID) (models.Job, error) {
	return a.jobs.Get(ctx, userID)
}
