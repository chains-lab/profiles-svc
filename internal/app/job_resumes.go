package app

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/app/domain"
	"github.com/google/uuid"
)

func (a App) UpdateDegree(ctx context.Context, userID uuid.UUID, degree string) error {
	return a.jobResumes.UpdateDegree(ctx, userID, degree)
}

func (a App) AdminUpdateDegree(ctx context.Context, userID uuid.UUID, degree string) error {
	return a.jobResumes.AdminUpdate(ctx, userID, entities.AdminJobUpdate{
		Degree: &degree,
	})
}

func (a App) UpdateIndustry(ctx context.Context, userID uuid.UUID, industry string) error {
	return a.jobResumes.UpdateIndustry(ctx, userID, industry)
}

func (a App) AdminUpdateIndustry(ctx context.Context, userID uuid.UUID, industry string) error {
	return a.jobResumes.AdminUpdate(ctx, userID, entities.AdminJobUpdate{
		Industry: &industry,
	})
}

func (a App) UpdateIncome(ctx context.Context, userID uuid.UUID, income string) error {
	return a.jobResumes.UpdateIncome(ctx, userID, income)
}

func (a App) AdminUpdateIncome(ctx context.Context, userID uuid.UUID, income string) error {
	return a.jobResumes.AdminUpdate(ctx, userID, entities.AdminJobUpdate{
		Income: &income,
	})
}
