package profiles

import (
	"context"

	svc "github.com/chains-lab/profiles-proto/gen/go/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problems"
	responses "github.com/chains-lab/profiles-svc/internal/api/grpc/response"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetOwnProfile(ctx context.Context, req *svc.GetOwnProfileRequest) (*svc.Profile, error) {
	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return nil, problems.UnauthenticatedError(ctx, "invalid initiator ID format")
	}

	profile, err := s.app.GetProfileByUserID(ctx, initiatorID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get profile by ID")

		return nil, err
	}

	return responses.Profile(profile), nil
}
