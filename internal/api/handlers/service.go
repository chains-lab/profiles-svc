package handlers

import (
	"context"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/api/interceptors"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/elector-cab-svc/internal/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
)

type App interface {
	UpdateResidence(ctx context.Context, userID uuid.UUID, city, country string) (models.Biography, error)
	UpdateSex(ctx context.Context, userID uuid.UUID, sex string) (models.Biography, error)
	UpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) (models.Biography, error)
	UpdateNationality(ctx context.Context, userID uuid.UUID, nationality string) (models.Biography, error)
	UpdatePrimaryLanguage(ctx context.Context, userID uuid.UUID, primaryLanguage string) (models.Biography, error)
	UpdateDegree(ctx context.Context, userID uuid.UUID, degree string) (models.JobResume, error)
	UpdateIndustry(ctx context.Context, userID uuid.UUID, industry string) (models.JobResume, error)
	UpdateIncome(ctx context.Context, userID uuid.UUID, income string) (models.JobResume, error)

	CreateCabinet(ctx context.Context, userID uuid.UUID) error
	GetCabinetByUserID(ctx context.Context, userID uuid.UUID) (models.Cabinet, error)
	GetCabinetByUsername(ctx context.Context, username string) (models.Cabinet, error)
	GetUserBiographyByUserID(ctx context.Context, userID uuid.UUID) (models.Biography, error)
	GetUserJobResumeByID(ctx context.Context, userID uuid.UUID) (models.JobResume, error)
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, profile app.UpdateProfileInput) (models.Profile, error)
}

type Service struct {
	app App
	cfg config.Config

	svc.UserServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}

func Log(ctx context.Context, requestID uuid.UUID) *logrus.Entry {
	entry, ok := ctx.Value(interceptors.LogCtxKey).(*logrus.Entry)
	if !ok {
		entry = logrus.NewEntry(logrus.New())
	}
	return entry.WithField("request_id", requestID)
}

func Meta(ctx context.Context) interceptors.MetaData {
	md, ok := ctx.Value(interceptors.MetaCtxKey).(interceptors.MetaData)
	if !ok {
		return interceptors.MetaData{}
	}
	return md
}
