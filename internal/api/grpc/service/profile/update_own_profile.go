package profile

import (
	"context"
	"time"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/profiles-proto/gen/go/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problem"
	responses "github.com/chains-lab/profiles-svc/internal/api/grpc/response"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateOwnProfile(ctx context.Context, req *svc.UpdateOwnProfileRequest) (*svc.Profile, error) {
	initiatorID, err := s.allowedRoles(ctx, req.Initiator, "update own profile", roles.User)
	input := app.UpdateProfileInput{
		Pseudonym:   req.Pseudonym,
		Description: req.Description,
		Avatar:      req.Avatar,
		Sex:         req.Sex,
	}

	if req.BirthDate != nil {
		birthdate, err := time.Parse(time.RFC3339, *req.BirthDate)
		if err != nil {
			logger.Log(ctx).WithError(err).Error("invalid birth date format")

			return nil, problem.InvalidArgumentError(ctx, "invalid ", &errdetails.BadRequest_FieldViolation{
				Field:       "birth_date",
				Description: "invalid date format, expected RFC3339",
			})
		}
		input.BirthDate = &birthdate
	}

	profile, err := s.app.UpdateProfile(ctx, initiatorID, input)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update user profile")

		return nil, err
	}

	logger.Log(ctx).Infof("profile for user %s has been updated", initiatorID)

	return responses.Profile(profile), nil
}
