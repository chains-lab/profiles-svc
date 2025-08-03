package service

import (
	"context"

	"github.com/chains-lab/profiles-svc/internal/api/responses"
	"github.com/chains-lab/profiles-svc/internal/logger"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/profile"
	"github.com/google/uuid"
)

func (s Service) UpdateOfficialByAdmin(ctx context.Context, req *svc.UpdateOfficialByAdminRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("invalid user ID format")

		return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	profile, err := s.app.AdminUpdateProfileOfficial(ctx, userID, req.Official)

	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to update field official in profile")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("official status for user %s has been updated by admin %s", userID, meta.InitiatorID)

	return responses.Profile(profile), nil
}
