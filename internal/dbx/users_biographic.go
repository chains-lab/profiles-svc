package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const UsersBiographic = "users_biographic"

type BiographicModel struct {
	UserID                   uuid.UUID  `db:"user_id"`
	Sex                      *string    `db:"sex,omitempty"`
	Birthday                 *time.Time `db:"birthday,omitempty"`
	Citizenship              *string    `db:"citizenship,omitempty"`
	Nationality              *string    `db:"nationality,omitempty"`
	PrimaryLanguage          *string    `db:"primary_language,omitempty"`
	SexUpdatedAt             *time.Time `db:"sex_updated_at,omitempty"`
	CitizenshipUpdatedAt     *time.Time `db:"citizenship_updated_at,omitempty"`
	NationalityUpdatedAt     *time.Time `db:"nationality_updated_at,omitempty"`
	PrimaryLanguageUpdatedAt *time.Time `db:"primary_language_updated_at,omitempty"`
}

type BiographicQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewBiographic(db *sql.DB) BiographicQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return BiographicQ{
		db:       db,
		selector: builder.Select("*").From(UsersBiographic),
		inserter: builder.Insert(UsersBiographic),
		updater:  builder.Update(UsersBiographic),
		deleter:  builder.Delete(UsersBiographic),
		counter:  builder.Select("COUNT(*) AS count").From(UsersBiographic),
	}
}

func (q BiographicQ) Insert(ctx context.Context, m BiographicModel) error {
	values := map[string]interface{}{
		"user_id":                     m.UserID,
		"sex":                         m.Sex,
		"birthday":                    m.Birthday,
		"citizenship":                 m.Citizenship,
		"nationality":                 m.Nationality,
		"primary_language":            m.PrimaryLanguage,
		"sex_updated_at":              m.SexUpdatedAt,
		"citizenship_updated_at":      m.CitizenshipUpdatedAt,
		"nationality_updated_at":      m.NationalityUpdatedAt,
		"primary_language_updated_at": m.PrimaryLanguageUpdatedAt,
	}
	// squirrel автоматически подставит NULL для nil
	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return err
	}
	_, err = q.db.ExecContext(ctx, query, args...)
	return err
}

type UpdatePersonalitiesInput struct {
	Sex                      *string
	SexUpdatedAt             *time.Time
	Citizenship              *string
	CitizenshipUpdatedAt     *time.Time
	Nationality              *string
	NationalityUpdatedAt     *time.Time
	PrimaryLanguage          *string
	PrimaryLanguageUpdatedAt *time.Time
}

func (q BiographicQ) Update(ctx context.Context, input UpdatePersonalitiesInput) error {
	updates := make(map[string]interface{})

	if input.Sex != nil {
		updates["sex"] = *input.Sex
	}
	if input.SexUpdatedAt != nil {
		updates["sex_updated_at"] = *input.SexUpdatedAt
	}
	if input.Citizenship != nil {
		updates["citizenship"] = *input.Citizenship
	}
	if input.CitizenshipUpdatedAt != nil {
		updates["citizenship_updated_at"] = *input.CitizenshipUpdatedAt
	}
	if input.Nationality != nil {
		updates["nationality"] = *input.Nationality
	}
	if input.NationalityUpdatedAt != nil {
		updates["nationality_updated_at"] = *input.NationalityUpdatedAt
	}
	if input.PrimaryLanguage != nil {
		updates["primary_language"] = *input.PrimaryLanguage
	}
	if input.PrimaryLanguageUpdatedAt != nil {
		updates["primary_language_updated_at"] = *input.PrimaryLanguageUpdatedAt
	}

	query, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return err
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

func (q BiographicQ) Select(ctx context.Context) ([]BiographicModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var personalities []BiographicModel
	for rows.Next() {
		var personality BiographicModel
		err := rows.Scan(
			&personality.UserID,
			&personality.Sex,
			&personality.Birthday,
			&personality.Citizenship,
			&personality.Nationality,
			&personality.PrimaryLanguage,
			&personality.SexUpdatedAt,
			&personality.CitizenshipUpdatedAt,
			&personality.NationalityUpdatedAt,
			&personality.PrimaryLanguageUpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		personalities = append(personalities, personality)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return personalities, nil
}

func (q BiographicQ) Get(ctx context.Context) (BiographicModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return BiographicModel{}, err
	}

	row := q.db.QueryRowContext(ctx, query, args...)
	var personality BiographicModel
	err = row.Scan(
		&personality.UserID,
		&personality.Sex,
		&personality.Birthday,
		&personality.Citizenship,
		&personality.Nationality,
		&personality.PrimaryLanguage,
		&personality.SexUpdatedAt,
		&personality.CitizenshipUpdatedAt,
		&personality.NationalityUpdatedAt,
		&personality.PrimaryLanguageUpdatedAt,
	)

	return personality, nil
}

func (q BiographicQ) FilterByUserID(userID uuid.UUID) BiographicQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	return q
}

func (q BiographicQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return err
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q BiographicQ) Count(ctx context.Context) (int64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, err
	}

	var count int64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q BiographicQ) Page(limit, offset uint64) BiographicQ {
	q.counter = q.counter.Limit(limit).Offset(offset)
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q BiographicQ) Transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
