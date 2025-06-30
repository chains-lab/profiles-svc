package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) UpdateOwnIncome(ctx context.Context, req *svc.UpdateOwnIncomeRequest) (*svc.JobResume, error) {
	meta := Meta(ctx)

	job, err := s.app.UpdateIncome(ctx, meta.InitiatorID, req.Income)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update user income")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.JobResume(job), nil
}
