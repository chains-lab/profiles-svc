package controller

import (
	"context"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/chains-lab/profiles-svc/internal/domain/modules/profile"
	"github.com/google/uuid"
)

type Domain interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, username string) (entity.Profile, error)

	FilterProfile(ctx context.Context, params profile.FilterParams, offset, limit int32) (entity.ProfileCollection, error)

	GetProfileByID(ctx context.Context, userID uuid.UUID) (entity.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (entity.Profile, error)

	UpdateProfile(ctx context.Context, accountID uuid.UUID, input profile.UpdateParams) (entity.Profile, error)
	UpdateProfileOfficial(ctx context.Context, accountID uuid.UUID, official bool) (entity.Profile, error)
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (entity.Profile, error)
}

type Service struct {
	domain Domain
	log    logium.Logger
}

func New(log logium.Logger, profile Domain) Service {
	return Service{
		domain: profile,
		log:    log,
	}
}
