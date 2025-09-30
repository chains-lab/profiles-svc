package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/profiles-svc/internal/data/pgdb"
	"github.com/chains-lab/profiles-svc/internal/domain/models"
	"github.com/chains-lab/profiles-svc/internal/domain/services/profile"
	"github.com/google/uuid"
)

func (d *Database) CreateProfile(ctx context.Context, profile models.Profile) error {
	schema := profileModelToSchema(profile)

	return d.sql.profiles.New().Insert(ctx, schema)
}

func (d *Database) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	row, err := d.sql.profiles.New().FilterUserID(userID).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Profile{}, nil
		}

		return models.Profile{}, err
	}

	return profileSchemaToModel(row), nil
}

func (d *Database) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	row, err := d.sql.profiles.New().FilterUsername(username).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Profile{}, nil
		}
		return models.Profile{}, err
	}

	return profileSchemaToModel(row), nil
}

func (d *Database) FilterProfiles(
	ctx context.Context,
	filters profile.FilterParams,
	page uint64,
	size uint64,
) (models.ProfileCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.profiles.New()

	if filters.Username != nil {
		query.FilterUsernameLike(*filters.Username)
	}
	if filters.Pseudonym != nil {
		query.FilterPseudonymLike(*filters.Pseudonym)
	}
	if filters.Official != nil {
		query.FilterOfficial(*filters.Official)
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.ProfileCollection{}, err
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.ProfileCollection{}, err
	}

	result := make([]models.Profile, len(rows))
	for i, profile := range rows {
		result[i] = profileSchemaToModel(profile)
	}

	return models.ProfileCollection{
		Data:  result,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d *Database) UpdateProfile(
	ctx context.Context,
	userID uuid.UUID,
	input profile.Update,
	updatedAt time.Time,
) error {
	q := d.sql.profiles.New().FilterUserID(userID)

	if input.Pseudonym != nil {
		if *input.Pseudonym == "" {
			q.UpdatePseudonym(nil)
		} else {
			q.UpdatePseudonym(input.Pseudonym)
		}
	}
	if input.Description != nil {
		if *input.Description == "" {
			q.UpdateDescription(nil)
		} else {
			q.UpdateDescription(input.Description)
		}
	}
	if input.Avatar != nil {
		if *input.Avatar == "" {
			q.UpdateAvatar(nil)
		} else {
			q.UpdateAvatar(input.Avatar)
		}
	}

	return q.Update(ctx, updatedAt)
}

func (d *Database) UpdateProfileUsername(
	ctx context.Context,
	userID uuid.UUID,
	username string,
	updatedAt time.Time,
) error {
	return d.sql.profiles.New().FilterUserID(userID).UpdateUsername(username).Update(ctx, updatedAt)
}

func (d *Database) UpdateProfileBirthDate(
	ctx context.Context,
	userID uuid.UUID,
	birthDate time.Time,
	updatedAt time.Time,
) error {
	return d.sql.profiles.New().FilterUserID(userID).UpdateBirthDate(birthDate).Update(ctx, updatedAt)
}

func (d *Database) UpdateProfileSex(
	ctx context.Context,
	userID uuid.UUID,
	sex string,
	updatedAt time.Time,
) error {
	return d.sql.profiles.New().FilterUserID(userID).UpdateSex(sex).Update(ctx, updatedAt)
}

func (d *Database) UpdateProfileOfficial(
	ctx context.Context,
	userID uuid.UUID,
	official bool,
	updatedAt time.Time,
) error {
	return d.sql.profiles.New().FilterUserID(userID).UpdateOfficial(official).Update(ctx, updatedAt)
}

func profileSchemaToModel(p pgdb.Profile) models.Profile {
	return models.Profile{
		UserID:      p.UserID,
		Username:    p.Username,
		Pseudonym:   p.Pseudonym,
		Description: p.Description,
		Avatar:      p.Avatar,
		Sex:         p.Sex,
		BirthDate:   p.BirthDate,
		Official:    p.Official,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func profileModelToSchema(p models.Profile) pgdb.Profile {
	return pgdb.Profile{
		UserID:      p.UserID,
		Username:    p.Username,
		Pseudonym:   p.Pseudonym,
		Description: p.Description,
		Avatar:      p.Avatar,
		Sex:         p.Sex,
		BirthDate:   p.BirthDate,
		Official:    p.Official,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
