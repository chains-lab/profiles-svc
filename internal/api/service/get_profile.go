package service

import (
	"context"

	"github.com/chains-lab/profile-svc/internal/api/responses"
	"github.com/chains-lab/profile-svc/internal/app/models"
	"github.com/chains-lab/profile-svc/internal/logger"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/profile"
	"github.com/google/uuid"
)

func (s Service) GetProfile(ctx context.Context, req *svc.GetProfileRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	var profile models.Profile
	var err error

	if req.GetUserId() != "" {
		userID, err := uuid.Parse(req.GetUserId())
		if err != nil {
			logger.Log(ctx, meta.RequestID).WithError(err).Error("invalid user ID format")

			return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
				Field:       "user_id",
				Description: "invalid UUID format for user ID",
			})
		}

		profile, err = s.app.GetProfileByUserID(ctx, userID)
		if err != nil {
			logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to get profile by user ID")

			return nil, responses.AppError(ctx, meta.RequestID, err)
		}
	} else if req.GetUsername() != "" {
		profile, err = s.app.GetProfileByUsername(ctx, req.GetUsername())
		if err != nil {
			logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to get profile by username")

			return nil, responses.AppError(ctx, meta.RequestID, err)
		}
	} else {
		return nil, responses.BadRequestError(ctx, meta.RequestID,
			responses.Violation{Field: "username", Description: "username is required"},
			responses.Violation{Field: "user_id", Description: "user_id is required"})
	}

	return responses.Profile(profile), nil
}
