package handlers

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
	"github.com/google/uuid"
)

func (s Service) GetOwnProfile(ctx context.Context, _ *svc.Empty) (*svc.Profile, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	profile, err := s.app.GetProfileByUserID(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to get profile by ID")

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.Profile(profile), nil
}
