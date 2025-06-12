package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/profile-storage/internal/config"
	"github.com/chains-lab/profile-storage/internal/repo/redisdb"
	"github.com/chains-lab/profile-storage/internal/repo/sqldb"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type ProfileModel struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Pseudonym   *string   `json:"pseudonym,omitempty"`
	Description *string   `json:"description,omitempty"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	Official    bool      `json:"official"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type profilesRedis interface {
	Set(ctx context.Context, profile redisdb.ProfileModel)
	GetByID(ctx context.Context, id uuid.UUID) (redisdb.ProfileModel, error)
	GetByUsername(ctx context.Context, username string) (redisdb.ProfileModel, error)
}

type profilesSql interface {
	New() sqldb.ProfilesQ
	Insert(ctx context.Context, input sqldb.ProfileInsertInput) error
	Update(ctx context.Context, input sqldb.UpdateProfileInput) error
	Select(ctx context.Context) ([]sqldb.ProfileModel, error)
	Get(ctx context.Context) (sqldb.ProfileModel, error)
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)

	FilterByID(id uuid.UUID) sqldb.ProfilesQ
	FilterByUsername(username string) sqldb.ProfilesQ

	Page(limit, offset uint64) sqldb.ProfilesQ

	Transaction(fn func(ctx context.Context) error) error
}

type ProfileRepo struct {
	sql   profilesSql
	redis profilesRedis
}

func NewProfileRepo(cfg config.Config, log *logrus.Logger) (ProfileRepo, error) {
	db, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return ProfileRepo{}, err
	}

	redisImpl := redisdb.NewProfiles(
		cfg.Database.Redis.Lifetime,
		log,
		redis.NewClient(&redis.Options{
			Addr:     cfg.Database.Redis.Addr,
			Password: cfg.Database.Redis.Password,
			DB:       cfg.Database.Redis.DB,
		}),
	)

	sqlImpl := sqldb.NewProfiles(db)

	return ProfileRepo{
		sql:   sqlImpl,
		redis: &redisImpl,
	}, nil
}

type ProfileCreateInput struct {
	ID        uuid.UUID
	Username  string
	Official  bool
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (p ProfileRepo) Create(ctx context.Context, input ProfileCreateInput) error {
	profile := redisdb.ProfileModel{
		ID:        input.ID,
		Username:  input.Username,
		Official:  input.Official,
		UpdatedAt: input.UpdatedAt,
		CreatedAt: input.CreatedAt,
	}

	p.redis.Set(ctx, profile)

	sqlInput := sqldb.ProfileInsertInput{
		ID:        input.ID,
		Username:  input.Username,
		Official:  input.Official,
		UpdatedAt: input.UpdatedAt,
		CreatedAt: input.CreatedAt,
	}

	return p.sql.New().Insert(ctx, sqlInput)
}

type ProflieUpdateInput struct {
	Username    *string
	Pseudonym   *string
	Description *string
	AvatarURL   *string
	Official    *bool
	UpdatedAt   time.Time
}

func (p ProfileRepo) Update(ctx context.Context, ID uuid.UUID, input ProflieUpdateInput) error {
	updates := sqldb.UpdateProfileInput{
		Username:    input.Username,
		Pseudonym:   input.Pseudonym,
		Description: input.Description,
		AvatarURL:   input.AvatarURL,
		Official:    input.Official,
	}

	if err := p.sql.New().FilterByID(ID).Update(ctx, updates); err != nil {
		return err
	}

	profile, err := p.sql.New().FilterByID(ID).Get(ctx)
	if err != nil {
		return err
	}

	redisInput := redisdb.ProfileModel{
		ID:          ID,
		Username:    profile.Username,
		Pseudonym:   profile.Pseudonym,
		Description: profile.Description,
		AvatarURL:   profile.AvatarURL,
		Official:    profile.Official,
		UpdatedAt:   profile.UpdatedAt,
	}

	p.redis.Set(ctx, redisInput)

	return nil
}

func (p ProfileRepo) GetByID(ctx context.Context, ID uuid.UUID) (ProfileModel, error) {
	redisProfile, err := p.redis.GetByID(ctx, ID)
	if !errors.Is(err, redis.Nil) {
		return ProfileModel{
			ID:          redisProfile.ID,
			Username:    redisProfile.Username,
			Pseudonym:   redisProfile.Pseudonym,
			Description: redisProfile.Description,
			AvatarURL:   redisProfile.AvatarURL,
			Official:    redisProfile.Official,
			UpdatedAt:   redisProfile.UpdatedAt,
			CreatedAt:   redisProfile.CreatedAt,
		}, nil
	}

	sqlProfile, err := p.sql.New().FilterByID(ID).Get(ctx)
	if err != nil {
		return ProfileModel{}, err
	}

	return ProfileModel{
		ID:          sqlProfile.ID,
		Username:    sqlProfile.Username,
		Pseudonym:   sqlProfile.Pseudonym,
		Description: sqlProfile.Description,
		AvatarURL:   sqlProfile.AvatarURL,
		Official:    sqlProfile.Official,
		UpdatedAt:   sqlProfile.UpdatedAt,
		CreatedAt:   sqlProfile.CreatedAt,
	}, nil
}

func (p ProfileRepo) GetByUsername(ctx context.Context, username string) (ProfileModel, error) {
	redisProfile, err := p.redis.GetByUsername(ctx, username)
	if !errors.Is(err, redis.Nil) {
		return ProfileModel{
			ID:          redisProfile.ID,
			Username:    redisProfile.Username,
			Pseudonym:   redisProfile.Pseudonym,
			Description: redisProfile.Description,
			AvatarURL:   redisProfile.AvatarURL,
			Official:    redisProfile.Official,
			UpdatedAt:   redisProfile.UpdatedAt,
			CreatedAt:   redisProfile.CreatedAt,
		}, nil
	}

	sqlProfile, err := p.sql.New().FilterByUsername(username).Get(ctx)
	if err != nil {
		return ProfileModel{}, err
	}

	return ProfileModel{
		ID:          sqlProfile.ID,
		Username:    sqlProfile.Username,
		Pseudonym:   sqlProfile.Pseudonym,
		Description: sqlProfile.Description,
		AvatarURL:   sqlProfile.AvatarURL,
		Official:    sqlProfile.Official,
		UpdatedAt:   sqlProfile.UpdatedAt,
		CreatedAt:   sqlProfile.CreatedAt,
	}, nil
}

func (p ProfileRepo) Transaction(fn func(ctx context.Context) error) error {
	return p.sql.New().Transaction(fn)
}
