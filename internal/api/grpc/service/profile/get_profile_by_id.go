package profile

import (
	"context"

	svc "github.com/chains-lab/profiles-proto/gen/go/svc/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problems"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/responses"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetProfileById(ctx context.Context, req *svc.GetProfileByIdRequest) (*svc.Profile, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return nil, problems.InvalidArgumentError(ctx, "invalid format user id", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	profile, err := s.app.GetProfileByUserID(ctx, userID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get profile by user ID")

		return nil, err
	}

	return responses.Profile(profile), nil
}
