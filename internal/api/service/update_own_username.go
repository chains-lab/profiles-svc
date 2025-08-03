package service

import (
	"context"

	"github.com/chains-lab/profiles-proto/gen/go/svc"
	"github.com/chains-lab/profiles-svc/internal/api/responses"
	"github.com/chains-lab/profiles-svc/internal/logger"
)

func (s Service) UpdateOwnUsername(ctx context.Context, req *svc.UpdateOwnUsernameRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	profile, err := s.app.UpdateUsername(ctx, meta.InitiatorID, req.Username)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to update user profile")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("username for user %s has been updated to %s", meta.InitiatorID, req.Username)

	return responses.Profile(profile), nil
}
