package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) UpdateOwnDegree(ctx context.Context, req *svc.UpdateOwnDegreeRequest) (*svc.JobResume, error) {
	meta := Meta(ctx)

	job, err := s.app.UpdateDegree(ctx, meta.InitiatorID, req.Degree)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update user degree")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.JobResume(job), nil
}
