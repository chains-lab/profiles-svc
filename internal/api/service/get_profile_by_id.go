package service

import (
	"context"

	"github.com/chains-lab/profiles-proto/gen/go/svc"
	"github.com/chains-lab/profiles-svc/internal/api/responses"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetProfileById(ctx context.Context, req *svc.GetProfileByIdRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("invalid user ID format")

		return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	profile, err := s.app.GetProfileByUserID(ctx, userID)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to get profile by user ID")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Profile(profile), nil
}
