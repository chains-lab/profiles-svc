package profile

import (
	"context"
	"time"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/profiles-proto/gen/go/svc/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/meta"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problems"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/responses"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/logger"
)

func (s Service) CreateOwnProfile(ctx context.Context, req *svc.CreateProfileRequest) (*svc.Profile, error) {
	user := meta.User(ctx)

	if user.Role != roles.User {
		logger.Log(ctx).Error("user does not have permission to create profile")

		return nil, problems.UnauthenticatedError(ctx, "user does not have permission to create profile")
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

	profile, err := s.app.CreateProfile(ctx, user.ID, input)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to create profile")

		return nil, err
	}

	logger.Log(ctx).Infof("created profile for user %s", user.ID)

	return responses.Profile(profile), nil
}
