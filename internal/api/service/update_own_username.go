package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) UpdateOwnUsername(ctx context.Context, req *svc.UpdateOwnUsernameRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	profile, err := s.app.UpdateProfile(ctx, meta.InitiatorID, app.UpdateProfileInput{
		Username: &req.Username,
	})
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update user profile")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Profile(profile), nil
}
