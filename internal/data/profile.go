package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/profiles-svc/internal/data/pgdb"
	"github.com/chains-lab/profiles-svc/internal/domain/models"
	"github.com/chains-lab/profiles-svc/internal/domain/services/profile"
	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
)

func (d *Database) CreateProfile(ctx context.Context, profile models.Profile) error {
	schema := profileModelToSchema(profile)

	return d.sql.profiles.New().Insert(ctx, schema)
}

func (d *Database) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	row, err := d.sql.profiles.New().FilterUserID(userID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Profile{}, nil
	case err != nil:
		return models.Profile{}, err
	}

	return profileSchemaToModel(row), nil
}

func (d *Database) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	row, err := d.sql.profiles.New().FilterUsername(username).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Profile{}, nil
	case err != nil:
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

	if filters.UserID != nil {
		query.FilterUserID(filters.UserID...)
	}
	if filters.UsernameLike != nil {
		query.FilterUsernameLike(*filters.UsernameLike)
	} else if len(filters.Username) > 0 {
		query.FilterUsername(filters.Username...)
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

	result := make([]models.Profile, 0, len(rows))
	for _, p := range rows {
		result = append(result, profileSchemaToModel(p))
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
	if input.Sex != nil {
		q.UpdateSex(*input.Sex)
	}
	if input.BirthDate != nil {
		q.UpdateBirthDate(*input.BirthDate)
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

func (d *Database) ResetProfile(ctx context.Context, userID uuid.UUID, username string, resetAt time.Time) error {
	err := d.sql.profiles.New().FilterUserID(userID).
		UpdateUsername(username).
		UpdatePseudonym(nil).
		UpdateDescription(nil).
		UpdateAvatar(nil).
		Update(ctx, resetAt)
	if err != nil {
		return err
	}

	return nil
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
