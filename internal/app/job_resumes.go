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

func (a App) UpdateIndustry(ctx context.Context, userID uuid.UUID, industry string) (models.JobResume, error) {
	err := a.jobResumes.UpdateIndustry(ctx, userID, industry)
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

type AdminUpdateJobResumeInput struct {
	Degree   *string
	Industry *string
	Income   *string
}

func (a App) AdminUpdateJobResume(ctx context.Context, userID uuid.UUID, input AdminUpdateJobResumeInput) (models.JobResume, error) {
	jobResume, err := a.jobResumes.Get(ctx, userID)
	if err != nil {
		return models.JobResume{}, err
	}

	err = a.jobResumes.AdminUpdate(ctx, userID, entities.AdminJobUpdate{
		Degree:   input.Degree,
		Industry: input.Industry,
		Income:   input.Income,
	})
	if err != nil {
		return models.JobResume{}, err
	}

	return jobResume, nil
}
