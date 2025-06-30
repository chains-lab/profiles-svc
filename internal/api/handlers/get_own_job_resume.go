package handlers

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
	"github.com/google/uuid"
)

func (s Service) GetOwnJobResume(ctx context.Context, _ *svc.Empty) (*svc.JobResume, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	job, err := s.app.GetUserJobResumeByID(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to get job")

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.JobResume(job), nil
}
