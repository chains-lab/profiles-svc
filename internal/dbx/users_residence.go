package dbx

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const UsersResidenceTable = "users_residence"

type UsersResidenceModel struct {
	UserID           uuid.UUID  `db:"user_id"`
	Country          *string    `db:"country,omitempty"`
	City             *string    `db:"city,omitempty"`
	CountryUpdatedAt *time.Time `db:"country_updated_at,omitempty"`
	CityUpdatedAt    *time.Time `db:"city_updated_at,omitempty"`
}

type ResidenceQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewResidence(db *sql.DB) ResidenceQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return ResidenceQ{
		db:       db,
		selector: builder.Select("*").From(UsersResidenceTable),
		inserter: builder.Insert(UsersResidenceTable),
		updater:  builder.Update(UsersResidenceTable),
		deleter:  builder.Delete(UsersResidenceTable),
		counter:  builder.Select("COUNT(*) AS count").From(UsersResidenceTable),
	}
}

func (q ResidenceQ) Insert(ctx context.Context, m UsersResidenceModel) error {
	values := map[string]interface{}{
		"user_id":            m.UserID,
		"country":            m.Country,
		"city":               m.City,
		"country_updated_at": m.CountryUpdatedAt,
		"city_updated_at":    m.CityUpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return err
	}

	_, err = q.db.ExecContext(ctx, query, args...)
	return err
}

type UpdateResidenceInput struct {
	Country          *string
	CountryUpdatedAt *time.Time
	City             *string
	CityUpdatedAt    *time.Time
}

func (q ResidenceQ) Update(ctx context.Context, userID uuid.UUID, input UpdateResidenceInput) error {
	values := map[string]interface{}{}
	if input.Country != nil {
		values["country"] = input.Country
	}
	if input.CountryUpdatedAt != nil {
		values["country_updated_at"] = input.CountryUpdatedAt
	}
	if input.City != nil {
		values["city"] = input.City
	}
	if input.CityUpdatedAt != nil {
		values["city_updated_at"] = input.CityUpdatedAt
	}

	query, args, err := q.updater.
		SetMap(values).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = q.db.ExecContext(ctx, query, args...)
	return err
}
