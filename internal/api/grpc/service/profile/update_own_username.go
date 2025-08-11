package profile

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/profiles-proto/gen/go/profile"
	responses "github.com/chains-lab/profiles-svc/internal/api/grpc/response"
	"github.com/chains-lab/profiles-svc/internal/logger"
)

func (s Service) UpdateOwnUsername(ctx context.Context, req *svc.UpdateOwnUsernameRequest) (*svc.Profile, error) {
	initiatorID, err := s.allowedRoles(ctx, req.Initiator, "reset profile", roles.User)

	profile, err := s.app.UpdateUsername(ctx, initiatorID, req.Username)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update user profile")

		return nil, err
	}

	logger.Log(ctx).Infof("username for user %s has been updated to %s", initiatorID, req.Username)

	return responses.Profile(profile), nil
}
