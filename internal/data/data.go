package data

import (
	"context"
	"database/sql"

	"github.com/chains-lab/profiles-svc/internal/data/pgdb"
)

type Database struct {
	sql SqlDB
}

func NewDatabase(db *sql.DB) *Database {
	profilesSql := pgdb.NewProfiles(db)

	return &Database{
		sql: SqlDB{
			profiles: profilesSql,
		},
	}
}
func (d *Database) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.sql.profiles.New().Transaction(ctx, fn)
}

type SqlDB struct {
	profiles pgdb.ProfilesQ
}
