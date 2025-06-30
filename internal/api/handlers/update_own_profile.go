package handlers

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
	"github.com/google/uuid"
)

func (s Service) UpdateOwnProfile(ctx context.Context, req *svc.UpdateOwnProfileRequest) (*svc.Profile, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	profile, err := s.app.UpdateProfile(ctx, meta.InitiatorID, app.UpdateProfileInput{
		Username:    req.Username,
		Pseudonym:   req.Pseudonym,
		Description: req.Description,
		Avatar:      req.Avatar,
	})
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to update user profile")

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.Profile(profile), nil
}
