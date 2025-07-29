package service

import (
	"context"
	"time"

	"github.com/chains-lab/profile-svc/internal/api/responses"
	"github.com/chains-lab/profile-svc/internal/app"
	"github.com/chains-lab/profile-svc/internal/logger"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/profile"
)

func (s Service) UpdateOwnProfile(ctx context.Context, req *svc.UpdateOwnProfileRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	input := app.UpdateProfileInput{
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

	profile, err := s.app.UpdateProfile(ctx, meta.InitiatorID, input)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to update user profile")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("profile for user %s has been updated", meta.InitiatorID)

	return responses.Profile(profile), nil
}
