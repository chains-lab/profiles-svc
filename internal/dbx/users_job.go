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

	UserDegree   *string `db:"user_degree,omitempty"`
	UserIndustry *string `db:"user_industry,omitempty"`
	UserIncome   *string `db:"user_income,omitempty"`

	UserDegreeUpdatedAt   *time.Time `db:"user_degree_updated_at,omitempty"`
	UserIndustryUpdatedAt *time.Time `db:"user_industry_updated_at,omitempty"`
	UserIncomeUpdatedAt   *time.Time `db:"user_income_updated_at,omitempty"`
}

type JobQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewJob(db *sql.DB) JobQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return JobQ{
		db:       db,
		selector: builder.Select("*").From(UsersJobTable),
		inserter: builder.Insert(UsersJobTable),
		updater:  builder.Update(UsersJobTable),
		deleter:  builder.Delete(UsersJobTable),
		counter:  builder.Select("COUNT(*) AS count").From(UsersJobTable),
	}
}

func (q JobQ) Insert(ctx context.Context, input JobModel) error {
	values := map[string]interface{}{
		"user_id":                  input.UserID,
		"user_degree":              input.UserDegree,
		"user_industry":            input.UserIndustry,
		"user_income":              input.UserIncome,
		"user_degree_updated_at":   input.UserDegreeUpdatedAt,
		"user_industry_updated_at": input.UserIndustryUpdatedAt,
		"user_income_updated_at":   input.UserIncomeUpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return err
	}

	_, err = q.db.ExecContext(ctx, query, args...)
	return err
}

type UpdateJobInput struct {
	UserDegree            *string
	UserDegreeUpdatedAt   *time.Time
	UserIndustry          *string
	UserIndustryUpdatedAt *time.Time
	UserIncome            *string
	UserIncomeUpdatedAt   *time.Time
}

func (q JobQ) Update(ctx context.Context, userID uuid.UUID, input UpdateJobInput) error {
	values := map[string]interface{}{}
	if input.UserDegree != nil {
		values["user_degree"] = input.UserDegree
	}
	if input.UserDegreeUpdatedAt != nil {
		values["user_degree_updated_at"] = input.UserDegreeUpdatedAt
	}
	if input.UserIndustry != nil {
		values["user_industry"] = input.UserIndustry
	}
	if input.UserIndustryUpdatedAt != nil {
		values["user_industry_updated_at"] = input.UserIndustryUpdatedAt
	}
	if input.UserIncome != nil {
		values["user_income"] = input.UserIncome
	}
	if input.UserIncomeUpdatedAt != nil {
		values["user_income_updated_at"] = input.UserIncomeUpdatedAt
	}

	query, args, err := q.updater.SetMap(values).Where(sq.Eq{"user_id": userID}).ToSql()
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

func (q JobQ) Select(ctx context.Context) ([]JobModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []JobModel
	for rows.Next() {
		var job JobModel
		if err := rows.Scan(
			&job.UserID,
			&job.UserDegree,
			&job.UserIndustry,
			&job.UserIncome,
			&job.UserDegreeUpdatedAt,
			&job.UserIndustryUpdatedAt,
			&job.UserIncomeUpdatedAt,
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

func (q JobQ) Get(ctx context.Context) (JobModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return JobModel{}, err
	}

	row := q.db.QueryRowContext(ctx, query, args...)
	var job JobModel
	err = row.Scan(
		&job.UserID,
		&job.UserDegree,
		&job.UserIndustry,
		&job.UserIncome,
		&job.UserDegreeUpdatedAt,
		&job.UserIndustryUpdatedAt,
		&job.UserIncomeUpdatedAt,
	)

	return job, err
}

func (q JobQ) Delete(ctx context.Context, userID uuid.UUID) error {
	query, args, err := q.deleter.Where(sq.Eq{"user_id": userID}).ToSql()
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

func (q JobQ) Count(ctx context.Context) (int64, error) {
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

func (q JobQ) FilterByUserID(userID uuid.UUID) JobQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	return q
}

func (q JobQ) Page(limit, offset uint64) JobQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(limit).Offset(offset)
	return q
}

func (q JobQ) Transaction(fn func(ctx context.Context) error) error {
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
