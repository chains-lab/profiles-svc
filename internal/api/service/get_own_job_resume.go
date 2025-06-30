package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetOwnJobResume(ctx context.Context, _ *emptypb.Empty) (*svc.JobResume, error) {
	meta := Meta(ctx)

	job, err := s.app.GetUserJobResumeByID(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to get job")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.JobResume(job), nil
}
