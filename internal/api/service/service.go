package service

import (
	"context"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/api/interceptors"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/elector-cab-svc/internal/config"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"github.com/google/uuid"
)

type App interface {
	CreateProfileAndCabinet(ctx context.Context, userID uuid.UUID, input app.CreateCabinetInput) (models.Profile, error)

	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)

	GetCabinetByUserID(ctx context.Context, userID uuid.UUID) (models.Cabinet, error)
	GetCabinetByUsername(ctx context.Context, username string) (models.Cabinet, error)

	GetBiographyByUserID(ctx context.Context, userID uuid.UUID) (models.Biography, error)

	UpdateProfile(ctx context.Context, userID uuid.UUID, profile app.UpdateProfileInput) (models.Profile, error)
	UpdateUsername(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error)

	UpdateSex(ctx context.Context, userID uuid.UUID, sex string) (models.Biography, error)
	UpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) (models.Biography, error)
	UpdateResidence(ctx context.Context, userID uuid.UUID, input app.UpdateResidenceInput) (models.Biography, error)

	AdminUpdateBiography(ctx context.Context, userID uuid.UUID, input app.UpdateBiographyInput) (models.Biography, error)
	AdminUpdateProfile(ctx context.Context, userID uuid.UUID, input app.AdminUpdateProfileInput) (models.Profile, error)

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
