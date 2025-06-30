package app

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/app/entities"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) UpdateDegree(ctx context.Context, userID uuid.UUID, degree string) (models.JobResume, error) {
	err := a.jobResumes.UpdateDegree(ctx, userID, degree)
	if err != nil {
		return models.JobResume{}, err
	}

	return models.JobResume{}, nil
}

func (a App) AdminUpdateDegree(ctx context.Context, userID uuid.UUID, degree string) (models.JobResume, error) {
	err := a.jobResumes.AdminUpdate(ctx, userID, entities.AdminJobUpdate{
		Degree: &degree,
	})
	if err != nil {
		return models.JobResume{}, err
	}

	return models.JobResume{}, nil
}

func (a App) UpdateIndustry(ctx context.Context, userID uuid.UUID, industry string) (models.JobResume, error) {
	err := a.jobResumes.UpdateIndustry(ctx, userID, industry)
	if err != nil {
		return models.JobResume{}, err
	}

	return models.JobResume{}, nil
}

func (a App) AdminUpdateIndustry(ctx context.Context, userID uuid.UUID, industry string) (models.JobResume, error) {
	err := a.jobResumes.AdminUpdate(ctx, userID, entities.AdminJobUpdate{
		Industry: &industry,
	})
	if err != nil {
		return models.JobResume{}, err
	}

	return models.JobResume{}, nil
}

func (a App) UpdateIncome(ctx context.Context, userID uuid.UUID, income string) (models.JobResume, error) {
	err := a.jobResumes.UpdateIncome(ctx, userID, income)
	if err != nil {
		return models.JobResume{}, err
	}

	return models.JobResume{}, nil
}

func (a App) AdminUpdateIncome(ctx context.Context, userID uuid.UUID, income string) (models.JobResume, error) {
	err := a.jobResumes.AdminUpdate(ctx, userID, entities.AdminJobUpdate{
		Income: &income,
	})
	if err != nil {
		return models.JobResume{}, err
	}

	return models.JobResume{}, nil
}
