package dbx

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const UsersJobResumeTable = "users_job_resume"

type JobResumeModel struct {
	UserID uuid.UUID `db:"user_id"`

	Degree   *string `db:"degree,omitempty"`
	Industry *string `db:"industry,omitempty"`
	Income   *string `db:"income,omitempty"`

	DegreeUpdatedAt   *time.Time `db:"degree_updated_at,omitempty"`
	IndustryUpdatedAt *time.Time `db:"industry_updated_at,omitempty"`
	IncomeUpdatedAt   *time.Time `db:"income_updated_at,omitempty"`
}

type JobResumesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewJobs(db *sql.DB) JobResumesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return JobResumesQ{
		db:       db,
		selector: builder.Select("*").From(UsersJobResumeTable),
		inserter: builder.Insert(UsersJobResumeTable),
		updater:  builder.Update(UsersJobResumeTable),
		deleter:  builder.Delete(UsersJobResumeTable),
		counter:  builder.Select("COUNT(*) AS count").From(UsersJobResumeTable),
	}
}

func (q JobResumesQ) New() JobResumesQ {
	return NewJobs(q.db)
}

func (q JobResumesQ) Insert(ctx context.Context, input JobResumeModel) error {
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

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
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

func (q JobResumesQ) Update(ctx context.Context, input UpdateJobInput) error {
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

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q JobResumesQ) Get(ctx context.Context) (JobResumeModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return JobResumeModel{}, err
	}

	var job JobResumeModel
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
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

func (q JobResumesQ) Select(ctx context.Context) ([]JobResumeModel, error) {
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

	var jobs []JobResumeModel
	for rows.Next() {
		var job JobResumeModel
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

func (q JobResumesQ) Delete(ctx context.Context) error {
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

func (q JobResumesQ) FilterUserID(userID uuid.UUID) JobResumesQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})

	return q
}

func (q JobResumesQ) Count(ctx context.Context) (int, error) {
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

func (q JobResumesQ) Page(limit, offset uint64) JobResumesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(limit).Offset(offset)

	return q
}
