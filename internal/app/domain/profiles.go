package domain

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/chains-lab/profile-storage/internal/app/ape"
	"github.com/chains-lab/profile-storage/internal/app/models"
	"github.com/chains-lab/profile-storage/internal/config"
	"github.com/chains-lab/profile-storage/internal/repo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type profilesRepo interface {
	Create(ctx context.Context, input repo.ProfileCreateInput) error
	Update(ctx context.Context, ID uuid.UUID, input repo.ProflieUpdateInput) error
	GetByID(ctx context.Context, ID uuid.UUID) (repo.ProfileModel, error)
	GetByUsername(ctx context.Context, username string) (repo.ProfileModel, error)
}

type Profiles struct {
	repo profilesRepo
}

func NewProfiles(cfg config.Config, log *logrus.Logger) (*Profiles, error) {
	data, err := repo.NewProfileRepo(cfg, log)
	if err != nil {
		return nil, err
	}

	return &Profiles{
		repo: data,
	}, nil
}

func (p *Profiles) Create(ctx context.Context, ID uuid.UUID, createdAt time.Time) *ape.Error {
	username, err := generateElectorID()
	if err != nil {
		return ape.ErrorInternal(err)
	}

	err = p.repo.Create(ctx, repo.ProfileCreateInput{
		ID:        ID,
		Username:  username,
		Official:  false,
		UpdatedAt: createdAt,
		CreatedAt: createdAt,
	})
	if err != nil {
		return ape.ErrorInternal(err)
	}

	return nil
}

type UpdateInput struct {
	Username    *string
	Pseudonym   *string
	Description *string
	AvatarURL   *string
	Official    *bool
}

func (p *Profiles) Update(ctx context.Context, ID uuid.UUID, input UpdateInput) *ape.Error {
	err := p.repo.Update(ctx, ID, repo.ProflieUpdateInput{
		Username:    input.Username,
		Pseudonym:   input.Pseudonym,
		Description: input.Description,
		AvatarURL:   input.AvatarURL,
		Official:    input.Official,
	})
	if err != nil {
		switch { //TODO
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

func (p *Profiles) GetByID(ctx context.Context, ID uuid.UUID) (models.Profile, *ape.Error) {
	profile, err := p.repo.GetByID(ctx, ID)
	if err != nil {
		switch {
		default: //TODO
			return models.Profile{}, ape.ErrorInternal(err)
		}
	}

	return parseRepoModel(profile), nil
}

func (p *Profiles) GetByUsername(ctx context.Context, username string) (models.Profile, *ape.Error) {
	profile, err := p.repo.GetByUsername(ctx, username)
	if err != nil {
		switch {
		default: //TODO
			return models.Profile{}, ape.ErrorInternal(err)
		}
	}

	return parseRepoModel(profile), nil
}

func parseRepoModel(data repo.ProfileModel) models.Profile {
	return models.Profile{
		ID:          data.ID,
		Username:    data.Username,
		Pseudonym:   data.Pseudonym,
		Description: data.Description,
		AvatarURL:   data.AvatarURL,
		Official:    data.Official,
		UpdatedAt:   data.UpdatedAt,
		CreatedAt:   data.CreatedAt,
	}
}

func generateElectorID() (string, error) {
	const length = 10
	const prefix = "elector"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		b := make([]byte, 1)
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}
		result[i] = '0' + (b[0] % 10)
	}

	return prefix + string(result), nil
}
