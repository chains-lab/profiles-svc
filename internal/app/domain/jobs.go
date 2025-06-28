package domain

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/elector-cab-svc/internal/dbx"
	"github.com/google/uuid"
)

type JobQ interface {
	New() dbx.JobsQ

	Insert(ctx context.Context, input dbx.JobModel) error
	Update(ctx context.Context, input dbx.UpdateJobInput) error
	Select(ctx context.Context) ([]dbx.JobModel, error)
	Get(ctx context.Context) (dbx.JobModel, error)
	Delete(ctx context.Context) error

	FilterUserID(userID uuid.UUID) dbx.JobsQ

	Page(limit, offset uint64) dbx.JobsQ
	Count(ctx context.Context) (int, error)
}

type Jobs struct {
	queries JobQ
}

func NewJob(db *sql.DB) (Jobs, error) {
	return Jobs{
		queries: dbx.NewJobs(db),
	}, nil
}

func (j Jobs) Create(ctx context.Context, userID uuid.UUID) error {
	if err := j.queries.Insert(ctx, dbx.JobModel{
		UserID: userID,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func (j Jobs) Get(ctx context.Context, userID uuid.UUID) (models.Job, error) {
	job, err := j.queries.FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		default:
			return models.Job{}, ape.ErrorInternal(err) //TODO
		}
	}

	return JobFromDb(job), nil
}

func (j Jobs) UpdateDegree(ctx context.Context, userID uuid.UUID, degree string) error {
	job, err := j.Get(ctx, userID)
	if err != nil {
		return ape.ErrorInternal(err) //TODO
	}

	if !(models.ValidateDegree(degree)) {
		return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
	}

	now := time.Now().UTC()

	if job.DegreeUpdatedAt != nil {
		last := *job.DegreeUpdatedAt

		if now.Sub(last) < 365*24*time.Hour {
			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}
	}

	if err := j.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateJobInput{
		Degree:          &degree,
		DegreeUpdatedAt: &now,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func (j Jobs) UpdateIndustry(ctx context.Context, userID uuid.UUID, industry string) error {
	job, err := j.Get(ctx, userID)
	if err != nil {
		return ape.ErrorInternal(err) //TODO
	}

	if !(models.ValidateIndustry(industry)) {
		return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
	}

	now := time.Now().UTC()

	if job.IndustryUpdatedAt != nil {
		last := *job.IndustryUpdatedAt

		if now.Sub(last) < 365*24*time.Hour {
			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}
	}

	if err := j.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateJobInput{
		Industry:          &industry,
		IndustryUpdatedAt: &now,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func (j Jobs) UpdateIncome(ctx context.Context, userID uuid.UUID, income string) error {
	job, err := j.Get(ctx, userID)
	if err != nil {
		return ape.ErrorInternal(err) //TODO
	}

	if !(models.ValidateIncome(income)) {
		return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
	}

	now := time.Now().UTC()

	if job.IncomeUpdatedAt != nil {
		last := *job.IncomeUpdatedAt

		if now.Sub(last) < 365*24*time.Hour {
			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}
	}

	if err := j.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateJobInput{
		Income:          &income,
		IncomeUpdatedAt: &now,
	}); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

type AdminJobUpdate struct {
	Degree   *string `json:"degree"`
	Industry *string `json:"industry"`
	Income   *string `json:"income"`
}

func (j Jobs) AdminUpdate(ctx context.Context, userID uuid.UUID, input AdminJobUpdate) error {
	_, err := j.Get(ctx, userID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()

	var dbInput dbx.UpdateJobInput

	if input.Degree != nil {
		if !models.ValidateDegree(*input.Degree) {
			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}
		dbInput.Degree = input.Degree
		dbInput.DegreeUpdatedAt = &now
	}

	if input.Industry != nil {
		if !models.ValidateIndustry(*input.Industry) {
			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}
		dbInput.Industry = input.Industry
		dbInput.IndustryUpdatedAt = &now

	}

	if input.Income != nil {
		if !models.ValidateIncome(*input.Income) {
			return ape.ErrorInternal(fmt.Errorf("")) //TODO: add error
		}
		dbInput.Income = input.Income
		dbInput.IncomeUpdatedAt = &now

	}

	if err := j.queries.New().FilterUserID(userID).Update(ctx, dbInput); err != nil {
		switch {
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func JobFromDb(model dbx.JobModel) models.Job {
	return models.Job{
		UserID:   model.UserID,
		Degree:   model.Degree,
		Industry: model.Industry,
		Income:   model.Income,

		DegreeUpdatedAt:   model.DegreeUpdatedAt,
		IndustryUpdatedAt: model.IndustryUpdatedAt,
		IncomeUpdatedAt:   model.IncomeUpdatedAt,
	}
}
