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

func (s Service) UpdateOwnUsername(ctx context.Context, req *svc.UpdateOwnUsernameRequest) (*svc.Profile, error) {
	user := meta.User(ctx)
	if user.Role != roles.User {
		logger.Log(ctx).Error("user does not have permission to update own username")

		return nil, problems.UnauthenticatedError(ctx, "user does not have permission to update own username")
	}

	profile, err := s.app.UpdateUsername(ctx, user.ID, req.Username)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update user profile")

		return nil, err
	}

	logger.Log(ctx).Infof("username for user %s has been updated to %s", user.ID, req.Username)

	return responses.Profile(profile), nil
}
