package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"github.com/google/uuid"
)

func (s Service) UpdateOfficialByAdmin(ctx context.Context, req *svc.UpdateOfficialByAdminRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("invalid user ID format")

		return nil, responses.BadRequestError(ctx, meta.RequestID, ape.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	profile, err := s.app.AdminUpdateProfile(ctx, userID, app.AdminUpdateProfileInput{
		Official: &req.Official,
	})

	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update field official in profile")
		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Profile(profile), nil
}
