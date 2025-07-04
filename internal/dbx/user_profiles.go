package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const profilesTable = "profiles"

type ProfileModel struct {
	UserID      uuid.UUID `db:"user_id"`
	Username    string    `db:"username"`
	Pseudonym   *string   `db:"pseudonym,omitempty"`
	Description *string   `db:"description,omitempty"`
	Avatar      *string   `db:"avatar,omitempty"`
	Official    bool      `db:"official"`

	UsernameUpdatedAt time.Time `db:"username_updated_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	CreatedAt         time.Time `db:"created_at"`
}

type ProfilesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewProfiles(db *sql.DB) ProfilesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return ProfilesQ{
		db:       db,
		selector: builder.Select("*").From(profilesTable),
		inserter: builder.Insert(profilesTable),
		updater:  builder.Update(profilesTable),
		deleter:  builder.Delete(profilesTable),
		counter:  builder.Select("COUNT(*) AS count").From(profilesTable),
	}
}

func (q ProfilesQ) New() ProfilesQ {
	return NewProfiles(q.db)
}

func (q ProfilesQ) Insert(ctx context.Context, input ProfileModel) error {
	values := map[string]interface{}{
		"user_id":             input.UserID,
		"username":            input.Username,
		"pseudonym":           input.Pseudonym,
		"description":         input.Description,
		"avatar":              input.Avatar,
		"official":            input.Official,
		"username_updated_at": input.UsernameUpdatedAt,
		"updated_at":          input.UpdatedAt,
		"created_at":          input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for profiles: %w", err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

type UpdateProfileInput struct {
	Username          *string
	Pseudonym         *string
	Description       *string
	Avatar            *string
	Official          *bool
	UsernameUpdatedAt *time.Time
	UpdatedAt         time.Time
}

func (q ProfilesQ) Update(ctx context.Context, input UpdateProfileInput) error {
	updates := make(map[string]interface{})

	if input.Username != nil {
		updates["username"] = *input.Username
	}
	if input.UsernameUpdatedAt != nil {
		updates["username_updated_at"] = *input.UsernameUpdatedAt
	}
	if input.Pseudonym != nil {
		updates["pseudonym"] = *input.Pseudonym
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Avatar != nil {
		updates["avatar"] = *input.Avatar
	}
	if input.Official != nil {
		updates["official"] = *input.Official
	}

	query, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return err
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q ProfilesQ) Get(ctx context.Context) (ProfileModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return ProfileModel{}, err
	}

	var profile ProfileModel
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&profile.UserID,
		&profile.Username,
		&profile.Pseudonym,
		&profile.Description,
		&profile.Avatar,
		&profile.Official,
		&profile.UsernameUpdatedAt,
		&profile.UpdatedAt,
		&profile.CreatedAt,
	)

	return profile, err
}

func (q ProfilesQ) Select(ctx context.Context) ([]ProfileModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for profiles: %w", err)
	}

	var rows *sql.Rows

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []ProfileModel
	for rows.Next() {
		var profile ProfileModel
		err := rows.Scan(
			&profile.UserID,
			&profile.Username,
			&profile.Pseudonym,
			&profile.Description,
			&profile.Avatar,
			&profile.Official,
			&profile.UsernameUpdatedAt,
			&profile.UpdatedAt,
			&profile.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (q ProfilesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for profiles: %w", err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q ProfilesQ) FilterUserID(userID uuid.UUID) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})

	return q
}

func (q ProfilesQ) FilterUsername(username string) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"username": username})
	q.counter = q.counter.Where(sq.Eq{"username": username})
	q.deleter = q.deleter.Where(sq.Eq{"username": username})
	q.updater = q.updater.Where(sq.Eq{"username": username})

	return q
}

func (q ProfilesQ) Count(ctx context.Context) (int, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for profiles: %w", err)
	}

	var count int64
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (q ProfilesQ) Page(limit, offset uint64) ProfilesQ {
	q.counter = q.counter.Limit(limit).Offset(offset)
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
