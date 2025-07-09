package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/logger"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
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
