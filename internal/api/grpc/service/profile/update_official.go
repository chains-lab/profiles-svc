package profile

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/profiles-proto/gen/go/svc/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/meta"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problems"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/responses"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateOfficial(ctx context.Context, req *svc.UpdateOfficialRequest) (*svc.Profile, error) {
	user := meta.User(ctx)

	if user.Role != roles.Admin && user.Role != roles.SuperUser {
		logger.Log(ctx).Error("user does not have permission to update official")

		return nil, problems.UnauthenticatedError(ctx, "user does not have permission to update official")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return nil, problems.InvalidArgumentError(ctx, "invalid user ID format", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	profile, err := s.app.AdminUpdateProfileOfficial(ctx, userID, req.Official)

	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update field official in profile")

		return nil, err
	}

	logger.Log(ctx).Infof("official status for user %s has been updated by admin %s", userID, user.ID)

	return responses.Profile(profile), nil
}
