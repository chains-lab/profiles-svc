package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/umisto/profiles-svc/internal/domain/entity"
	"github.com/umisto/profiles-svc/internal/domain/modules/profile"
	"github.com/umisto/profiles-svc/internal/repo/pgdb"
)

func (r *Repository) CreateProfile(ctx context.Context, userID uuid.UUID, username string) (entity.Profile, error) {
	res, err := r.sql.profiles.New().Insert(ctx, pgdb.Profile{
		AccountID: userID,
		Username:  username,
	})
	if err != nil {
		return entity.Profile{}, err
	}

	return res.ToEntity(), nil
}

func (r *Repository) GetProfileByAccountID(ctx context.Context, accountId uuid.UUID) (entity.Profile, error) {
	row, err := r.sql.profiles.New().FilterAccountID(accountId).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Profile{}, nil
	case err != nil:
		return entity.Profile{}, err
	}

	return row.ToEntity(), nil
}

func (r *Repository) GetProfileByUsername(ctx context.Context, username string) (entity.Profile, error) {
	row, err := r.sql.profiles.New().FilterUsername(username).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Profile{}, nil
	case err != nil:
		return entity.Profile{}, err
	}

	return row.ToEntity(), nil
}

func (r *Repository) UpdateProfile(
	ctx context.Context,
	accountID uuid.UUID,
	input profile.UpdateParams,
) (entity.Profile, error) {
	q := r.sql.profiles.New().FilterAccountID(accountID)

	if input.Pseudonym != nil {
		q = q.UpdatePseudonym(input.Pseudonym)
	}
	if input.Description != nil {
		q = q.UpdateDescription(input.Description)
	}
	if input.Avatar != nil {
		q = q.UpdateAvatar(input.Avatar)
	}

	res, err := q.UpdateOne(ctx)
	if err != nil {
		return entity.Profile{}, err
	}

	return res.ToEntity(), nil
}

func (r *Repository) UpdateProfileUsername(
	ctx context.Context,
	accountID uuid.UUID,
	username string,
) (entity.Profile, error) {
	res, err := r.sql.profiles.New().
		FilterAccountID(accountID).
		UpdateUsername(username).
		UpdateOne(ctx)
	if err != nil {
		return entity.Profile{}, err
	}

	return res.ToEntity(), nil
}

func (r *Repository) UpdateProfileOfficial(
	ctx context.Context,
	accountID uuid.UUID,
	official bool,
) (entity.Profile, error) {
	res, err := r.sql.profiles.New().
		FilterAccountID(accountID).
		UpdateOfficial(official).
		UpdateOne(ctx)
	if err != nil {
		return entity.Profile{}, err
	}

	return res.ToEntity(), nil
}

func (r *Repository) FilterProfilesByUsername(
	ctx context.Context,
	prefix string,
	offset uint,
	limit uint,
) (entity.ProfileCollection, error) {
	rows, err := r.sql.profiles.New().
		FilterLikeUsername(prefix).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return entity.ProfileCollection{}, err
	}

	collection := make([]entity.Profile, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToEntity())
	}

	return entity.ProfileCollection{
		Data: collection,
		Page: uint(offset/limit) + 1,
		Size: uint(len(collection)),
	}, nil
}

func (r *Repository) FilterProfiles(
	ctx context.Context,
	params profile.FilterParams,
	offset uint,
	limit uint,
) (entity.ProfileCollection, error) {
	q := r.sql.profiles.New()

	if params.PseudonymPrefix != nil {
		q = q.FilterLikePseudonym(*params.PseudonymPrefix)
	}
	if params.UsernamePrefix != nil {
		q = q.FilterLikeUsername(*params.UsernamePrefix)
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return entity.ProfileCollection{}, err
	}

	collection := make([]entity.Profile, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToEntity())
	}

	return entity.ProfileCollection{
		Data: collection,
		Page: uint(offset/limit) + 1,
		Size: uint(len(collection)),
	}, nil
}

func (r *Repository) DeleteProfile(ctx context.Context, accountID uuid.UUID) error {
	return r.sql.profiles.New().FilterAccountID(accountID).Delete(ctx)
}
