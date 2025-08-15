package profile

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/profiles-proto/gen/go/svc/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/meta"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problems"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/responses"
	"github.com/chains-lab/profiles-svc/internal/logger"
)

func (s Service) GetOwnProfile(ctx context.Context, req *svc.GetOwnProfileRequest) (*svc.Profile, error) {
	user := meta.User(ctx)

	if user.Role != roles.User {
		logger.Log(ctx).Error("user does not have permission to get own profile")

		return nil, problems.UnauthenticatedError(ctx, "user does not have permission to get own profile")
	}

	profile, err := s.app.GetProfileByUserID(ctx, user.ID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get profile by ID")

		return nil, err
	}

	return responses.Profile(profile), nil
}
