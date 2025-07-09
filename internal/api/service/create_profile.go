package service

import (
	"context"
	"fmt"

	"github.com/chains-lab/elector-cab-svc/internal/ape"
	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	"github.com/chains-lab/elector-cab-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) CreateOwnProfile(ctx context.Context, req *svc.CreateProfileRequest) (*svc.Profile, error) {
	meta := Meta(ctx)

	if meta.Role != roles.User {
		logger.Log(ctx, meta.RequestID).Warnf(fmt.Sprintf(
			"user %s with role %s tried to create a cabinet, but only users can create profiles and cabinets",
			meta.InitiatorID, meta.Role),
		)

		return nil, responses.AppError(ctx, meta.RequestID, ape.RaiseOnlyUserCanHaveCabinetAndProfile(
			fmt.Errorf("user %s with role %s tried to create a profile and cabinet", meta.InitiatorID, meta.Role)),
		)
	}

	profile, err := s.app.CreateProfileAndCabinet(ctx, meta.InitiatorID, app.CreateCabinetInput{
		Username:    req.Username,
		Pseudonym:   req.Pseudonym,
		Description: req.Description,
		Avatar:      req.Avatar,
	})
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to create profile and cabinet")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("created profile and cabinet for user %s", meta.InitiatorID)

	return responses.Profile(profile), nil
}
