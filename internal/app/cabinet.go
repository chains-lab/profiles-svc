package app

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) CreateCabinet(ctx context.Context, userID uuid.UUID) (models.Cabinet, error) {
	txErr := a.transaction(func(ctx context.Context) error {
		err := a.profiles.Create(ctx, userID)
		if err != nil {
			return err
		}

		err = a.jobResumes.Create(ctx, userID)
		if err != nil {
			return err
		}

		err = a.biographies.Create(ctx, userID)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.Cabinet{}, txErr
	}

	return a.GetCabinetByUserID(ctx, userID)
}

func (a App) GetCabinetByUserID(ctx context.Context, userID uuid.UUID) (models.Cabinet, error) {
	profile, err := a.profiles.GetByID(ctx, userID)
	if err != nil {
		return models.Cabinet{}, err
	}

	jobs, err := a.jobResumes.Get(ctx, userID)
	if err != nil {
		return models.Cabinet{}, err
	}

	biography, err := a.biographies.GetByUserID(ctx, userID)
	if err != nil {
		return models.Cabinet{}, err
	}

	return models.Cabinet{
		Profile:   profile,
		Job:       jobs,
		Biography: biography,
	}, nil
}

func (a App) GetCabinetByUsername(ctx context.Context, username string) (models.Cabinet, error) {
	profile, err := a.profiles.GetByUsername(ctx, username)
	if err != nil {
		return models.Cabinet{}, err
	}

	jobs, err := a.jobResumes.Get(ctx, profile.UserID)
	if err != nil {
		return models.Cabinet{}, err
	}

	biography, err := a.biographies.GetByUserID(ctx, profile.UserID)
	if err != nil {
		return models.Cabinet{}, err
	}

	return models.Cabinet{
		Profile:   profile,
		Job:       jobs,
		Biography: biography,
	}, nil
}

func (a App) GetUserBiographyByUserID(ctx context.Context, userID uuid.UUID) (models.Biography, error) {
	return a.biographies.GetByUserID(ctx, userID)
}

func (a App) GetUserJobResumeByID(ctx context.Context, userID uuid.UUID) (models.JobResume, error) {
	return a.jobResumes.Get(ctx, userID)
}
