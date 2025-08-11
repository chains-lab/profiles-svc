package profile

import (
	"context"

	svc "github.com/chains-lab/profiles-proto/gen/go/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problem"
	responses "github.com/chains-lab/profiles-svc/internal/api/grpc/response"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) UpdateOwnUsername(ctx context.Context, req *svc.UpdateOwnUsernameRequest) (*svc.Profile, error) {
	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return nil, problem.UnauthenticatedError(ctx, "invalid initiator ID format")
	}

	profile, err := s.app.UpdateUsername(ctx, initiatorID, req.Username)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update user profile")

		return nil, err
	}

	logger.Log(ctx).Infof("username for user %s has been updated to %s", initiatorID, req.Username)

	return responses.Profile(profile), nil
}
