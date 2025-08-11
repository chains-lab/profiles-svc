package profile

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/profiles-proto/gen/go/profile"
	responses "github.com/chains-lab/profiles-svc/internal/api/grpc/response"
	"github.com/chains-lab/profiles-svc/internal/logger"
)

func (s Service) GetOwnProfile(ctx context.Context, req *svc.GetOwnProfileRequest) (*svc.Profile, error) {
	initiatorID, err := s.allowedRoles(ctx, req.Initiator, "gtt own profile", roles.User)
	if err != nil {
		return nil, err
	}

	profile, err := s.app.GetProfileByUserID(ctx, initiatorID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get profile by ID")

		return nil, err
	}

	return responses.Profile(profile), nil
}
