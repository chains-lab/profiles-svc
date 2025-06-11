package sqldb

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
	ID          uuid.UUID `db:"id"`
	Username    string    `db:"username"`
	Pseudonym   *string   `db:"pseudonym,omitempty"`
	Description *string   `db:"description,omitempty"`
	AvatarURL   *string   `db:"avatar_url,omitempty"`
	Official    bool      `db:"official"`
	UpdatedAt   time.Time `db:"updated_at"`
	CreatedAt   time.Time `db:"created_at"`
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

type ProfileInsertInput struct {
	ID        uuid.UUID
	Username  string
	CreatedAt time.Time
}

func (q ProfilesQ) Insert(ctx context.Context, input ProfileInsertInput) error {
	values := map[string]interface{}{
		"id":         input.ID,
		"username":   input.Username,
		"updated_at": input.CreatedAt,
		"created_at": input.CreatedAt,
	}

	query, args, err := q.inserter.
		Columns("id", "username", "created_at").
		Values(values["id"], values["username"], values["created_at"]).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

type UpdateProfileInput struct {
	Username    *string
	Pseudonym   *string
	Description *string
	AvatarURL   *string
	Official    *bool
}

func (q ProfilesQ) Update(ctx context.Context, input UpdateProfileInput) error {
	updates := make(map[string]interface{})

	if input.Username != nil {
		updates["username"] = *input.Username
	}
	if input.Pseudonym != nil {
		updates["pseudonym"] = *input.Pseudonym
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.AvatarURL != nil {
		updates["avatar_url"] = *input.AvatarURL
	}
	if input.Official != nil {
		updates["official"] = *input.Official
	}

	query, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for profiles: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (q ProfilesQ) Select(ctx context.Context) ([]ProfileModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for profiles: %w", err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []ProfileModel
	for rows.Next() {
		var profile ProfileModel
		err := rows.Scan(
			&profile.ID,
			&profile.Username,
			&profile.Pseudonym,
			&profile.Description,
			&profile.AvatarURL,
			&profile.Official,
			&profile.UpdatedAt,
			&profile.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning profile row: %w", err)
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (q ProfilesQ) Get(ctx context.Context) (ProfileModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return ProfileModel{}, fmt.Errorf("building get query for profiles: %w", err)
	}

	row := q.db.QueryRowContext(ctx, query, args...)
	var profile ProfileModel
	err = row.Scan(
		&profile.ID,
		&profile.Username,
		&profile.Pseudonym,
		&profile.Description,
		&profile.AvatarURL,
		&profile.Official,
		&profile.UpdatedAt,
		&profile.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ProfileModel{}, nil
		}
		return ProfileModel{}, fmt.Errorf("scanning profile row: %w", err)
	}

	return profile, nil
}

func (q ProfilesQ) FilterByID(id uuid.UUID) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	return q
}

func (q ProfilesQ) FilterByUsername(username string) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"username": username})
	q.counter = q.counter.Where(sq.Eq{"username": username})
	q.deleter = q.deleter.Where(sq.Eq{"username": username})
	q.updater = q.updater.Where(sq.Eq{"username": username})
	return q
}

func (q ProfilesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for profiles: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (q ProfilesQ) Count(ctx context.Context) (int64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for profiles: %w", err)
	}

	var count int64
	err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q ProfilesQ) Page(limit, offset uint64) ProfilesQ {
	q.counter = q.counter.Limit(limit).Offset(offset)
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
