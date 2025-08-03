package service

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/profiles-svc/internal/ape"
	"github.com/chains-lab/profiles-svc/internal/api/responses"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/logger"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/profile"
)

func (s Service) CreateOwnProfile(ctx context.Context, req *svc.CreateProfilrRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	if meta.Role != roles.User {
		logger.Log(ctx, meta.RequestID).Warnf(fmt.Sprintf(
			"user %s with role %s tried to create a profile, but only users can create profile",
			meta.InitiatorID, meta.Role),
		)

		return nil, responses.AppError(ctx, meta.RequestID, ape.RaiseOnlyUserCanHaveProfile(
			fmt.Errorf("user %s with role %s tried to create a profile", meta.InitiatorID, meta.Role)),
		)
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
			logger.Log(ctx, meta.RequestID).WithError(err).Error("invalid birth date format")

			return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
				Field:       "birth_date",
				Description: "invalid date format, expected RFC3339",
			})
		}
		input.BirthDate = &birthdate
	}

	profile, err := s.app.CreateProfile(ctx, meta.InitiatorID, input)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to create profile")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("created profile for user %s", meta.InitiatorID)

	return responses.Profile(profile), nil
}
