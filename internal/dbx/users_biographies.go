package dbx

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const UsersBiographies = "users_biographies"

type BioModel struct {
	UserID          uuid.UUID  `db:"user_id"`
	Sex             *string    `db:"sex,omitempty"`
	Birthday        *time.Time `db:"birthday,omitempty"`
	Nationality     *string    `db:"nationality,omitempty"`
	PrimaryLanguage *string    `db:"primary_language,omitempty"`
	Country         *string    `db:"country,omitempty"`
	City            *string    `db:"city,omitempty"`

	SexUpdatedAt             *time.Time `db:"sex_updated_at,omitempty"`
	NationalityUpdatedAt     *time.Time `db:"nationality_updated_at,omitempty"`
	PrimaryLanguageUpdatedAt *time.Time `db:"primary_language_updated_at,omitempty"`
	ResidenceUpdatedAt       *time.Time `db:"residence_updated_at,omitempty"`
}

type BiographiesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewBiographies(db *sql.DB) BiographiesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return BiographiesQ{
		db:       db,
		selector: builder.Select("*").From(UsersBiographies),
		inserter: builder.Insert(UsersBiographies),
		updater:  builder.Update(UsersBiographies),
		deleter:  builder.Delete(UsersBiographies),
		counter:  builder.Select("COUNT(*) AS count").From(UsersBiographies),
	}
}

func (q BiographiesQ) New() BiographiesQ {
	return NewBiographies(q.db)
}

func (q BiographiesQ) Insert(ctx context.Context, m BioModel) error {
	values := map[string]interface{}{
		"user_id":                     m.UserID,
		"sex":                         m.Sex,
		"birthday":                    m.Birthday,
		"nationality":                 m.Nationality,
		"primary_language":            m.PrimaryLanguage,
		"country":                     m.Country,
		"city":                        m.City,
		"sex_updated_at":              m.SexUpdatedAt,
		"nationality_updated_at":      m.NationalityUpdatedAt,
		"primary_language_updated_at": m.PrimaryLanguageUpdatedAt,
		"residence_updated_at":        m.ResidenceUpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
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

type UpdateBioInput struct {
	Birthday                 *time.Time
	Sex                      *string
	SexUpdatedAt             *time.Time
	Nationality              *string
	NationalityUpdatedAt     *time.Time
	PrimaryLanguage          *string
	PrimaryLanguageUpdatedAt *time.Time
	Country                  *string
	City                     *string
	ResidenceUpdatedAt       *time.Time
}

func (q BiographiesQ) Update(ctx context.Context, input UpdateBioInput) error {
	updates := make(map[string]interface{})

	if input.Birthday != nil {
		updates["birthday"] = *input.Birthday
	}
	if input.Sex != nil {
		updates["sex"] = *input.Sex
	}
	if input.SexUpdatedAt != nil {
		updates["sex_updated_at"] = *input.SexUpdatedAt
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
	if input.Country != nil {
		updates["country"] = *input.Country
	}
	if input.City != nil {
		updates["city"] = *input.City
	}
	if input.ResidenceUpdatedAt != nil {
		updates["residence_updated_at"] = *input.ResidenceUpdatedAt
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

func (q BiographiesQ) Get(ctx context.Context) (BioModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return BioModel{}, err
	}

	var personality BioModel
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&personality.UserID,
		&personality.Sex,
		&personality.Birthday,
		&personality.Nationality,
		&personality.PrimaryLanguage,
		&personality.Country,
		&personality.City,
		&personality.SexUpdatedAt,
		&personality.NationalityUpdatedAt,
		&personality.PrimaryLanguageUpdatedAt,
		&personality.ResidenceUpdatedAt,
	)

	return personality, nil
}

func (q BiographiesQ) Select(ctx context.Context) ([]BioModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, err
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

	var personalities []BioModel
	for rows.Next() {
		var personality BioModel
		err := rows.Scan(
			&personality.UserID,
			&personality.Sex,
			&personality.Birthday,
			&personality.Nationality,
			&personality.PrimaryLanguage,
			&personality.Country,
			&personality.City,
			&personality.SexUpdatedAt,
			&personality.NationalityUpdatedAt,
			&personality.PrimaryLanguageUpdatedAt,
			&personality.ResidenceUpdatedAt,
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

func (q BiographiesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
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

func (q BiographiesQ) FilterUserID(userID uuid.UUID) BiographiesQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})

	return q
}

func (q BiographiesQ) Count(ctx context.Context) (int, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, err
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

func (q BiographiesQ) Page(limit, offset uint64) BiographiesQ {
	q.counter = q.counter.Limit(limit).Offset(offset)
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}
