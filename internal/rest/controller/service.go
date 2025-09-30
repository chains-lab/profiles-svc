package controller

import (
	"context"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/domain/models"
	"github.com/chains-lab/profiles-svc/internal/domain/services/profile"
	"github.com/google/uuid"
)

type ProfileSvc interface {
	Create(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error)

	GetByID(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	GetByUsername(ctx context.Context, username string) (models.Profile, error)

	Filter(ctx context.Context, filters profile.FilterParams, page uint64, size uint64) (models.ProfileCollection, error)

	Update(ctx context.Context, userID uuid.UUID, input profile.Update) (models.Profile, error)
	UpdateUsername(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error)
	UpdateOfficial(ctx context.Context, userID uuid.UUID, official bool) (models.Profile, error)
	UpdateBirthDate(ctx context.Context, userID uuid.UUID, birthDate time.Time) (models.Profile, error)
	UpdateSex(ctx context.Context, userID uuid.UUID, sex string) (models.Profile, error)

	ResetUsername(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	ResetUserProfile(ctx context.Context, userID uuid.UUID) (models.Profile, error)
}

type Service struct {
	domain domain
	log    logium.Logger
}

type domain struct {
	Profile ProfileSvc
}

func New(log logium.Logger, profile ProfileSvc) Service {
	return Service{
		domain: domain{
			Profile: profile,
		},
		log: log,
	}
}
