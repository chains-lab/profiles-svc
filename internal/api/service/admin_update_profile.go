package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"github.com/google/uuid"
)

func (s Service) AdminUpdateProfile(ctx context.Context, req *svc.AdminUpdateProfileRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("invalid user ID format")

		return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	profile, err := s.app.AdminUpdateProfile(ctx, userID, app.AdminUpdateProfileInput{
		Username:    req.Username,
		Pseudonym:   req.Pseudonym,
		Description: req.Description,
		Avatar:      req.Avatar,
		Official:    req.Official,
	})

	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update profile")
		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Profile(profile), nil
}
