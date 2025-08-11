package profile

import (
	"context"
	"fmt"

	"github.com/chains-lab/profiles-proto/gen/go/common/userdata"
	profilesProto "github.com/chains-lab/profiles-proto/gen/go/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/problem"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/chains-lab/profiles-svc/internal/config"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"github.com/google/uuid"
)

type App interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, input app.CreateProfileInput) (models.Profile, error)

	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)

	UpdateProfile(ctx context.Context, userID uuid.UUID, profile app.UpdateProfileInput) (models.Profile, error)
	UpdateUsername(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error)
	AdminUpdateProfileOfficial(ctx context.Context, userID uuid.UUID, official bool) (models.Profile, error)

	ResetUsername(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	ResetUserProfile(ctx context.Context, userID uuid.UUID, input app.ResetUserProfileInput) (models.Profile, error)
}

type Service struct {
	app App
	cfg config.Config

	profilesProto.UnimplementedProfilesServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}

func (s Service) allowedRoles(ctx context.Context, req *userdata.UserData, action string, allowed ...string) (uuid.UUID, error) {
	initiatorID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return uuid.Nil, problem.UnauthenticatedError(ctx, "invalid initiator ID format")
	}

	allow := false
	for _, role := range allowed {
		if req.Role == role {
			allow = true
			break
		}
	}

	if !allow {
		logger.Log(ctx).Warnf(
			"user %s with role %s tried to perform this action: '%s', that requires one of the allowed roles: %v",
			req.UserId, req.Role, action, allowed,
		)

		return uuid.Nil, problem.PermissionDeniedError(ctx,
			fmt.Sprintf("initiator role can perform this '%s' action", action))
	}
	return initiatorID, nil
}
