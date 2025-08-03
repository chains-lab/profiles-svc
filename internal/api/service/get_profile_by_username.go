package service

import (
	"context"

	"github.com/chains-lab/profiles-proto/gen/go/svc"
	"github.com/chains-lab/profiles-svc/internal/api/responses"
	"github.com/chains-lab/profiles-svc/internal/logger"
)

func (s Service) GetProfileByUsername(ctx context.Context, req *svc.GetProfileByUsernameRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	profile, err := s.app.GetProfileByUsername(ctx, req.Username)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to get profile by username")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Profile(profile), nil
}
