package profile

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/profiles-proto/gen/go/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problem"
	responses "github.com/chains-lab/profiles-svc/internal/api/grpc/response"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) ResetProfileByAdmin(ctx context.Context, req *svc.ResetProfileByAdminRequest) (*svc.Profile, error) {
	if req.Initiator.Role == roles.Admin || req.Initiator.Role == roles.SuperUser {
		logger.Log(ctx).Error("unauthorized access")

		return nil, problem.PermissionDeniedError(ctx, "only admins roles can reset profile")
	}

	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return nil, problem.UnauthenticatedError(ctx, "invalid initiator ID format")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return nil, problem.InvalidArgumentError(ctx, "invalid user id format", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	profile, err := s.app.ResetUserProfile(ctx, userID, app.ResetUserProfileInput{
		Pseudonym:   req.Pseudonym,
		Description: req.Description,
		Avatar:      req.Avatar,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to reset profile")

		return nil, err
	}

	logger.Log(ctx).Infof("profile for user %s has been reset by admin %s", userID, initiatorID)

	return responses.Profile(profile), nil
}
