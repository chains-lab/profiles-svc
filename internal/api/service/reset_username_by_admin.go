package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/logger"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"github.com/google/uuid"
)

func (s Service) ResetUsernameByAdmin(ctx context.Context, req *svc.ResetUsernameByAdminRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("invalid user ID format")

		return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	profile, err := s.app.ResetUsername(ctx, userID)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to reset username")
		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("username for user %s has been reset by admin %s", userID, meta.InitiatorID)

	return responses.Profile(profile), nil
}
