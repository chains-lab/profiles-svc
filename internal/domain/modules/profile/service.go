package profile

import (
	"context"

	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/google/uuid"
)

type Service struct {
	db database
}

func New(db database) Service {
	return Service{
		db: db,
	}
}

type database interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, username string) (entity.Profile, error)

	GetProfileByAccountID(ctx context.Context, userID uuid.UUID) (entity.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (entity.Profile, error)

	UpdateProfile(ctx context.Context, userID uuid.UUID, params UpdateParams) (entity.Profile, error)

	UpdateProfileUsername(ctx context.Context, userID uuid.UUID, username string) (entity.Profile, error)
	UpdateProfileOfficial(ctx context.Context, userID uuid.UUID, official bool) (entity.Profile, error)

	DeleteProfile(ctx context.Context, userID uuid.UUID) error

	FilterProfiles(
		ctx context.Context,
		params FilterParams,
		offset uint,
		limit uint,
	) (entity.ProfileCollection, error)
}
