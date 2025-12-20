package repo

import (
	"database/sql"

	"github.com/umisto/profiles-svc/internal/repo/pgdb"
)

type Repository struct {
	sql SqlDB
}

type SqlDB struct {
	profiles pgdb.ProfilesQ
}

func New(db *sql.DB) *Repository {
	return &Repository{
		sql: SqlDB{
			profiles: pgdb.NewProfilesQ(db),
		},
	}
}
