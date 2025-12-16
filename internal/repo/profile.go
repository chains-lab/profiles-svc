package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/chains-lab/profiles-svc/internal/domain/modules/profile"
	"github.com/chains-lab/profiles-svc/internal/repo/pgdb"
	"github.com/google/uuid"
)

func (r *Repository) CreateProfile(ctx context.Context, userID uuid.UUID, username string) (entity.Profile, error) {
	row, err := r.sql.CreateProfile(ctx, pgdb.CreateProfileParams{
		AccountID: userID,
		Username:  username,
	})
	if err != nil {
		return entity.Profile{}, err
	}

	return row.ToEntity(), nil
}

func (r *Repository) GetProfileByUserID(ctx context.Context, accountId uuid.UUID) (entity.Profile, error) {
	row, err := r.sql.GetProfileByAccountID(ctx, accountId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Profile{}, nil
	case err != nil:
		return entity.Profile{}, err
	}

	return row.ToEntity(), nil
}

func (r *Repository) GetProfileByUsername(ctx context.Context, username string) (entity.Profile, error) {
	row, err := r.sql.GetProfileByUsername(ctx, username)
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
	row, err := r.sql.UpdateProfile(ctx, pgdb.UpdateProfileParams{
		AccountID: accountID,
		Pseudonym: sql.NullString{
			String: *input.Pseudonym,
			Valid:  input.Pseudonym == nil,
		},
		Description: sql.NullString{
			String: *input.Description,
			Valid:  input.Description == nil,
		},
		Avatar: sql.NullString{
			String: *input.Avatar,
			Valid:  input.Avatar == nil,
		},
	})
	if err != nil {
		return entity.Profile{}, err
	}

	return row.ToEntity(), nil
}

func (r *Repository) UpdateProfileUsername(
	ctx context.Context,
	userID uuid.UUID,
	username string,
) (entity.Profile, error) {
	row, err := r.sql.UpdateProfileUsername(ctx, pgdb.UpdateProfileUsernameParams{
		AccountID: userID,
		Username:  username,
	})
	if err != nil {
		return entity.Profile{}, err
	}

	return row.ToEntity(), nil
}

func (r *Repository) UpdateProfileOfficial(
	ctx context.Context,
	userID uuid.UUID,
	official bool,
) (entity.Profile, error) {
	row, err := r.sql.UpdateProfileOfficial(ctx, pgdb.UpdateProfileOfficialParams{
		AccountID: userID,
		Official:  official,
	})
	if err != nil {
		return entity.Profile{}, err
	}

	return row.ToEntity(), nil
}

func (r *Repository) FilterProfilesByUsername(
	ctx context.Context,
	prefix string,
	offset int32,
	limit int32,
) (entity.ProfileCollection, error) {
	rows, err := r.sql.ListProfilesByUsername(ctx, pgdb.ListProfilesByUsernameParams{
		Prefix: prefix,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return entity.ProfileCollection{}, err
	}

	collection := make([]entity.Profile, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToEntity())
	}

	return entity.ProfileCollection{
		Data: collection,
		Page: int32(offset/limit) + 1,
		Size: int32(len(collection)),
	}, nil
}

func (r *Repository) FilterProfilesByPseudonym(
	ctx context.Context,
	prefix string,
	offset int32,
	limit int32,
) (entity.ProfileCollection, error) {
	rows, err := r.sql.ListProfilesByPseudonym(ctx, pgdb.ListProfilesByPseudonymParams{
		Prefix: prefix,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return entity.ProfileCollection{}, err
	}

	collection := make([]entity.Profile, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToEntity())
	}

	return entity.ProfileCollection{
		Data: collection,
		Page: int32(offset/limit) + 1,
		Size: int32(len(collection)),
	}, nil
}

func (r *Repository) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return r.sql.DeleteProfile(ctx, userID)
}
