package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const UsersJobTable = "users_job"

type JobModel struct {
	UserID uuid.UUID `db:"user_id"`

	Degree   *string `db:"degree,omitempty"`
	Industry *string `db:"industry,omitempty"`
	Income   *string `db:"income,omitempty"`

	DegreeUpdatedAt   *time.Time `db:"degree_updated_at,omitempty"`
	IndustryUpdatedAt *time.Time `db:"industry_updated_at,omitempty"`
	IncomeUpdatedAt   *time.Time `db:"income_updated_at,omitempty"`
}

type JobsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewJobs(db *sql.DB) JobsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return JobsQ{
		db:       db,
		selector: builder.Select("*").From(UsersJobTable),
		inserter: builder.Insert(UsersJobTable),
		updater:  builder.Update(UsersJobTable),
		deleter:  builder.Delete(UsersJobTable),
		counter:  builder.Select("COUNT(*) AS count").From(UsersJobTable),
	}
}

func (q JobsQ) New() JobsQ {
	return NewJobs(q.db)
}

func (q JobsQ) Insert(ctx context.Context, input JobModel) error {
	values := map[string]interface{}{
		"user_id":             input.UserID,
		"degree":              input.Degree,
		"industry":            input.Industry,
		"income":              input.Income,
		"degree_updated_at":   input.DegreeUpdatedAt,
		"industry_updated_at": input.IndustryUpdatedAt,
		"income_updated_at":   input.IncomeUpdatedAt,
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

type UpdateJobInput struct {
	Degree            *string
	DegreeUpdatedAt   *time.Time
	Industry          *string
	IndustryUpdatedAt *time.Time
	Income            *string
	IncomeUpdatedAt   *time.Time
}

func (q JobsQ) Update(ctx context.Context, input UpdateJobInput) error {
	updates := map[string]interface{}{}

	if input.Degree != nil {
		updates["degree"] = input.Degree
	}
	if input.DegreeUpdatedAt != nil {
		updates["degree_updated_at"] = input.DegreeUpdatedAt
	}
	if input.Industry != nil {
		updates["industry"] = input.Industry
	}
	if input.IndustryUpdatedAt != nil {
		updates["industry_updated_at"] = input.IndustryUpdatedAt
	}
	if input.Income != nil {
		updates["income"] = input.Income
	}
	if input.IncomeUpdatedAt != nil {
		updates["income_updated_at"] = input.IncomeUpdatedAt
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

	return err
}

func (q JobsQ) Get(ctx context.Context) (JobModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return JobModel{}, err
	}

	var job JobModel
	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&job.UserID,
		&job.Degree,
		&job.Industry,
		&job.Income,
		&job.DegreeUpdatedAt,
		&job.IndustryUpdatedAt,
		&job.IncomeUpdatedAt,
	)

	return job, err
}

func (q JobsQ) Select(ctx context.Context) ([]JobModel, error) {
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

	var jobs []JobModel
	for rows.Next() {
		var job JobModel
		if err := rows.Scan(
			&job.UserID,
			&job.Degree,
			&job.Industry,
			&job.Income,
			&job.DegreeUpdatedAt,
			&job.IndustryUpdatedAt,
			&job.IncomeUpdatedAt,
		); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (q JobsQ) Delete(ctx context.Context) error {
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

func (q JobsQ) FilterUserID(userID uuid.UUID) JobsQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})

	return q
}

func (q JobsQ) Count(ctx context.Context) (int, error) {
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

func (q JobsQ) Page(limit, offset uint64) JobsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(limit).Offset(offset)

	return q
}

func (q JobsQ) Transaction(fn func(ctx context.Context) error) error {
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
