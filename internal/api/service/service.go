package service

import (
	"context"

	"github.com/chains-lab/profiles-svc/internal/api/interceptors"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/chains-lab/profiles-svc/internal/config"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/profile"
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

	svc.UserServiceServer
	svc.AdminServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}

func Meta(ctx context.Context) interceptors.MetaData {
	md, ok := ctx.Value(interceptors.MetaCtxKey).(interceptors.MetaData)
	if !ok {
		return interceptors.MetaData{}
	}
	return md
}
