package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const UsersResidenceTable = "users_residence"

type ResidenceModel struct {
	UserID    uuid.UUID  `db:"user_id"`
	Country   *string    `db:"country,omitempty"`
	City      *string    `db:"city,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type ResidencesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewResidences(db *sql.DB) ResidencesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return ResidencesQ{
		db:       db,
		selector: builder.Select("*").From(UsersResidenceTable),
		inserter: builder.Insert(UsersResidenceTable),
		updater:  builder.Update(UsersResidenceTable),
		deleter:  builder.Delete(UsersResidenceTable),
		counter:  builder.Select("COUNT(*) AS count").From(UsersResidenceTable),
	}
}

func (q ResidencesQ) New() ResidencesQ {
	return NewResidences(q.db)
}

func (q ResidencesQ) Insert(ctx context.Context, m ResidenceModel) error {
	values := map[string]interface{}{
		"user_id":    m.UserID,
		"country":    m.Country,
		"city":       m.City,
		"updated_at": m.UpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
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

type UpdateResidenceInput struct {
	Country   *string
	City      *string
	UpdatedAt *time.Time
}

func (q ResidencesQ) Update(ctx context.Context, input UpdateResidenceInput) error {
	values := map[string]interface{}{}
	if input.Country != nil {
		values["country"] = input.Country
	}
	if input.City != nil {
		values["city"] = input.City
	}
	if input.UpdatedAt != nil {
		values["updated_at"] = input.UpdatedAt
	}

	query, args, err := q.updater.SetMap(values).ToSql()
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

func (q ResidencesQ) Get(ctx context.Context) (ResidenceModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return ResidenceModel{}, err
	}

	var residence ResidenceModel
	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&residence.UserID,
		&residence.Country,
		&residence.City,
		&residence.UpdatedAt,
	)
	if err != nil {
		return ResidenceModel{}, err
	}

	return residence, nil
}

func (q ResidencesQ) Select(ctx context.Context) ([]ResidenceModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var residences []ResidenceModel
	for rows.Next() {
		var residence ResidenceModel
		err := rows.Scan(
			&residence.UserID,
			&residence.Country,
			&residence.City,
			&residence.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		residences = append(residences, residence)
	}

	return residences, nil
}

func (q ResidencesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
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

func (q ResidencesQ) FilterUserID(userID uuid.UUID) ResidencesQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	return q
}

func (q ResidencesQ) Count(ctx context.Context) (int, error) {
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

	return int(count), nil
}

func (q ResidencesQ) Page(limit, offset uint64) ResidencesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(limit).Offset(offset)

	return q
}

func (q ResidencesQ) Transaction(fn func(ctx context.Context) error) error {
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
