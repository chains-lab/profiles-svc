package handlers

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
	"github.com/google/uuid"
)

func (s Service) UpdateOwnIncome(ctx context.Context, req *svc.UpdateOwnIncomeRequest) (*svc.JobResume, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	job, err := s.app.UpdateIncome(ctx, meta.InitiatorID, req.Income)
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to update user income")

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.JobResume(job), nil
}
