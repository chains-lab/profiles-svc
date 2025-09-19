package profiles

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/chains-lab/profiles-svc/internal/dbx"
	"github.com/google/uuid"
)

type ProfileQ interface {
	New() dbx.ProfilesQ

	Insert(ctx context.Context, input dbx.ProfileModel) error
	Update(ctx context.Context, input map[string]any) error
	Get(ctx context.Context) (dbx.ProfileModel, error)
	Select(ctx context.Context) ([]dbx.ProfileModel, error)
	Delete(ctx context.Context) error

	FilterUserID(userID uuid.UUID) dbx.ProfilesQ
	FilterUsername(username string) dbx.ProfilesQ

	Count(ctx context.Context) (int, error)
	Page(limit, offset uint64) dbx.ProfilesQ
}

type Profiles struct {
	queries ProfileQ

	usernameRegex *regexp.Regexp
}

func NewProfile(db *sql.DB) (Profiles, error) {
	return Profiles{
		queries:       dbx.NewProfiles(db),
		usernameRegex: regexp.MustCompile(`^[a-zA-Z0-9._]{3,32}$`),
	}, nil
}

func ProfileFromDb(input dbx.ProfileModel) models.Profile {
	return models.Profile{
		UserID:      input.UserID,
		Username:    input.Username,
		Pseudonym:   input.Pseudonym,
		Description: input.Description,
		Avatar:      input.Avatar,
		Official:    input.Official,
		Sex:         input.Sex,
		Birthdate:   input.BirthDate,
		UpdatedAt:   input.UpdatedAt,
		CreatedAt:   input.CreatedAt,
	}
}
