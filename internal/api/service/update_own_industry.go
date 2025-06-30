package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) UpdateOwnIndustry(ctx context.Context, req *svc.UpdateOwnIndustryRequest) (*svc.JobResume, error) {
	meta := Meta(ctx)

	job, err := s.app.UpdateIndustry(ctx, meta.InitiatorID, req.Industry)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update user industry")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.JobResume(job), nil
}
