package profile

import (
	"context"
	"time"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/profiles-proto/gen/go/svc/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/meta"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problems"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/responses"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateOwnProfile(ctx context.Context, req *svc.UpdateOwnProfileRequest) (*svc.Profile, error) {
	user := meta.User(ctx)
	if user.Role != roles.User {
		logger.Log(ctx).Error("user does not have permission to update profile")

		return nil, problems.UnauthenticatedError(ctx, "user does not have permission to update profile")
	}

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

			return nil, problems.InvalidArgumentError(ctx, "invalid ", &errdetails.BadRequest_FieldViolation{
				Field:       "birth_date",
				Description: "invalid date format, expected RFC3339",
			})
		}
		input.BirthDate = &birthdate
	}

	profile, err := s.app.UpdateProfile(ctx, user.ID, input)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update user profile")

		return nil, err
	}

	logger.Log(ctx).Infof("profile for user %s has been updated", user.ID)

	return responses.Profile(profile), nil
}
