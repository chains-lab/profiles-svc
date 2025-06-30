package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetOwnProfile(ctx context.Context, _ *emptypb.Empty) (*svc.Profile, error) {
	meta := Meta(ctx)

	profile, err := s.app.GetProfileByUserID(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to get profile by ID")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Profile(profile), nil
}
