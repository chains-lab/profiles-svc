package dbx

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const BiographiesTable = "biographies"

type BioModel struct {
	UserID   uuid.UUID  `db:"user_id"`
	Birthday *time.Time `db:"birthday,omitempty"`
	Sex      *string    `db:"sex,omitempty"`
	City     *string    `db:"city,omitempty"`
	Region   *string    `db:"region,omitempty"`
	Country  *string    `db:"country,omitempty"`

	SexUpdatedAt       *time.Time `db:"sex_updated_at,omitempty"`
	ResidenceUpdatedAt *time.Time `db:"residence_updated_at,omitempty"`
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
		selector: builder.Select("*").From(BiographiesTable),
		inserter: builder.Insert(BiographiesTable),
		updater:  builder.Update(BiographiesTable),
		deleter:  builder.Delete(BiographiesTable),
		counter:  builder.Select("COUNT(*) AS count").From(BiographiesTable),
	}
}

func (q BiographiesQ) New() BiographiesQ {
	return NewBiographies(q.db)
}

func (q BiographiesQ) Insert(ctx context.Context, m BioModel) error {
	values := map[string]interface{}{
		"user_id":              m.UserID,
		"birthday":             m.Birthday,
		"sex":                  m.Sex,
		"city":                 m.City,
		"region":               m.Region,
		"country":              m.Country,
		"sex_updated_at":       m.SexUpdatedAt,
		"residence_updated_at": m.ResidenceUpdatedAt,
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
	Birthday           *time.Time
	Sex                *string
	SexUpdatedAt       *time.Time
	Country            *string
	Region             *string
	City               *string
	ResidenceUpdatedAt *time.Time
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
	if input.Country != nil {
		updates["country"] = *input.Country
	}
	if input.Region != nil {
		updates["region"] = *input.Region
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
		&personality.Birthday,
		&personality.Sex,
		&personality.City,
		&personality.Region,
		&personality.Country,
		&personality.SexUpdatedAt,
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
			&personality.Birthday,
			&personality.Sex,
			&personality.City,
			&personality.Region,
			&personality.Country,
			&personality.SexUpdatedAt,
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
