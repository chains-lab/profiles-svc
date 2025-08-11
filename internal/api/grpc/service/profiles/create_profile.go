package profiles

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/profiles-proto/gen/go/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problems"
	responses "github.com/chains-lab/profiles-svc/internal/api/grpc/response"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/logger"
)

func (s Service) CreateOwnProfile(ctx context.Context, req *svc.CreateProfileRequest) (*svc.Profile, error) {
	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return nil, problems.UnauthenticatedError(ctx, "invalid initiator ID format")
	}

	if req.Initiator.Role != roles.User {
		logger.Log(ctx).Warnf(fmt.Sprintf(
			"user %s with role %s tried to create a profile, but only users can create profile",
			req.Initiator.UserId, req.Initiator.Role),
		)

		return nil, problems.PermissionDeniedError(ctx, "only common users can create profile")
	}

	input := app.CreateProfileInput{
		Username:    req.Username,
		Pseudonym:   req.Pseudonym,
		Description: req.Description,
		Avatar:      req.Avatar,
		Sex:         req.Sex,
	}
	if req.BirthDate != nil {
		birthdate, err := time.Parse(time.RFC3339, *req.BirthDate)
		if err != nil {
			logger.Log(ctx).WithError(err).Error("invalid birth date format")

			return nil, problems.InvalidArgumentError(ctx, "birthdate is invalid format", &errdetails.BadRequest_FieldViolation{
				Field:       "birth_date",
				Description: "invalid date format, expected RFC3339",
			})
		}
		input.BirthDate = &birthdate
	}

	profile, err := s.app.CreateProfile(ctx, initiatorID, input)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to create profile")

		return nil, err
	}

	logger.Log(ctx).Infof("created profile for user %s", initiatorID)

	return responses.Profile(profile), nil
}
