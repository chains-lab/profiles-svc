package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const profilesTable = "profiles"

type Profile struct {
	UserID      uuid.UUID  `db:"user_id"`
	Username    string     `db:"username"`
	Official    bool       `db:"official"`
	Pseudonym   *string    `db:"pseudonym,omitempty"`
	Description *string    `db:"description,omitempty"`
	Avatar      *string    `db:"avatar,omitempty"`
	Sex         *string    `db:"sex"`
	BirthDate   *time.Time `db:"birth_date,omitempty"`

	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
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

func (q ProfilesQ) Insert(ctx context.Context, input Profile) error {
	values := map[string]interface{}{
		"user_id":     input.UserID,
		"username":    input.Username,
		"official":    input.Official,
		"pseudonym":   input.Pseudonym,
		"description": input.Description,
		"avatar":      input.Avatar,
		"sex":         input.Sex,
		"birth_date":  input.BirthDate,
		"updated_at":  input.UpdatedAt,
		"created_at":  input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for profile: %w", err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q ProfilesQ) Update(ctx context.Context, updatedAt time.Time) error {
	q.updater = q.updater.Set("updated_at", updatedAt)

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", profilesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q ProfilesQ) UpdateUsername(username string) ProfilesQ {
	q.updater = q.updater.Set("username", username)
	return q
}

func (q ProfilesQ) UpdatePseudonym(pseudonym *string) ProfilesQ {
	q.updater = q.updater.Set("pseudonym", pseudonym)
	return q
}

func (q ProfilesQ) UpdateDescription(description *string) ProfilesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q ProfilesQ) UpdateAvatar(avatar *string) ProfilesQ {
	q.updater = q.updater.Set("avatar", avatar)
	return q
}

func (q ProfilesQ) UpdateOfficial(official bool) ProfilesQ {
	q.updater = q.updater.Set("official", official)
	return q
}

func (q ProfilesQ) UpdateSex(sex string) ProfilesQ {
	q.updater = q.updater.Set("sex", sex)
	return q
}

func (q ProfilesQ) UpdateBirthDate(birthDate time.Time) ProfilesQ {
	q.updater = q.updater.Set("birth_date", birthDate)
	return q
}

func (q ProfilesQ) Get(ctx context.Context) (Profile, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Profile{}, err
	}

	var profile Profile
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&profile.UserID,
		&profile.Username,
		&profile.Official,
		&profile.Pseudonym,
		&profile.Description,
		&profile.Avatar,
		&profile.Sex,
		&profile.BirthDate,
		&profile.UpdatedAt,
		&profile.CreatedAt,
	)

	return profile, err
}

func (q ProfilesQ) Select(ctx context.Context) ([]Profile, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for profile: %w", err)
	}

	var rows *sql.Rows

	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var profile Profile
		err := rows.Scan(
			&profile.UserID,
			&profile.Username,
			&profile.Official,
			&profile.Pseudonym,
			&profile.Description,
			&profile.Avatar,
			&profile.Sex,
			&profile.BirthDate,
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
		return fmt.Errorf("building delete query for profile: %w", err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q ProfilesQ) FilterUserID(userID ...uuid.UUID) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})

	return q
}

func (q ProfilesQ) FilterUsername(username ...string) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"username": username})
	q.counter = q.counter.Where(sq.Eq{"username": username})
	q.deleter = q.deleter.Where(sq.Eq{"username": username})
	q.updater = q.updater.Where(sq.Eq{"username": username})

	return q
}

func (q ProfilesQ) FilterUsernameLike(username string) ProfilesQ {
	like := fmt.Sprintf("%%%s%%", username)
	q.selector = q.selector.Where(sq.Like{"username": like})
	q.counter = q.counter.Where(sq.Like{"username": like})
	return q
}

func (q ProfilesQ) FilterPseudonymLike(pseudonym string) ProfilesQ {
	like := fmt.Sprintf("%%%s%%", pseudonym)
	q.selector = q.selector.Where(sq.Like{"pseudonym": like})
	q.counter = q.counter.Where(sq.Like{"pseudonym": like})
	return q
}

func (q ProfilesQ) FilterOfficial(official bool) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"official": official})
	q.counter = q.counter.Where(sq.Eq{"official": official})
	return q
}

func (q ProfilesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for profile: %w", err)
	}

	var count uint64
	if tx, ok := TxFromCtx(ctx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
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

func (q ProfilesQ) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_, ok := TxFromCtx(ctx)
	if ok {
		return fn(ctx)
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
				err = fmt.Errorf("tx err: %v; rollback err: %v", err, rbErr)
			}
		}
	}()

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	if err = fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
