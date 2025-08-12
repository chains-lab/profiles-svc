package profileadmin

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/profiles-proto/gen/go/svc/profile"
	profileAdmiProto "github.com/chains-lab/profiles-proto/gen/go/svc/profileadmin"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/guard"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problem"
	responses "github.com/chains-lab/profiles-svc/internal/api/grpc/response"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateOfficial(ctx context.Context, req *profileAdmiProto.UpdateOfficialRequest) (*svc.Profile, error) {
	initiatorID, err := guard.AllowedRoles(ctx, req.Initiator, "update official status by admin",
		roles.Moder, roles.Admin, roles.SuperUser)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return nil, problem.InvalidArgumentError(ctx, "invalid user ID format", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	profile, err := s.app.AdminUpdateProfileOfficial(ctx, userID, req.Official)

	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update field official in profile")

		return nil, err
	}

	logger.Log(ctx).Infof("official status for user %s has been updated by admin %s", userID, initiatorID)

	return responses.Profile(profile), nil
}
