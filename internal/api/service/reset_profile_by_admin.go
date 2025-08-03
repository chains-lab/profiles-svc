package service

import (
	"context"

	"github.com/chains-lab/profiles-svc/internal/api/responses"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/logger"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/profile"
	"github.com/google/uuid"
)

func (s Service) ResetProfileByAdmin(ctx context.Context, req *svc.ResetProfileByAdminRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("invalid user ID format")

		return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	profile, err := s.app.ResetUserProfile(ctx, userID, app.ResetUserProfileInput{
		Pseudonym:   *req.Pseudonym,
		Description: *req.Description,
		Avatar:      *req.Avatar,
	})
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to reset profile")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("profile for user %s has been reset by admin %s", userID, meta.InitiatorID)

	return responses.Profile(profile), nil
}
