package profiles

import (
	"context"

	svc "github.com/chains-lab/profiles-proto/gen/go/profile"
	responses "github.com/chains-lab/profiles-svc/internal/api/grpc/response"
	"github.com/chains-lab/profiles-svc/internal/logger"
)

func (s Service) GetProfileByUsername(ctx context.Context, req *svc.GetProfileByUsernameRequest) (*svc.Profile, error) {
	profile, err := s.app.GetProfileByUsername(ctx, req.Username)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get profile by username")

		return nil, err
	}

	return responses.Profile(profile), nil
}
