package repo

import (
	"database/sql"

	"github.com/chains-lab/profiles-svc/internal/repo/pgdb"
)

type Repository struct {
	sql SqlDB
}

type SqlDB struct {
	inbox    pgdb.InboxEventsQ
	outbox   pgdb.OutboxEventsQ
	profiles pgdb.ProfilesQ
}

func New(db *sql.DB) *Repository {
	return &Repository{
		sql: SqlDB{
			inbox:    pgdb.NewInboxEventsQ(db),
			outbox:   pgdb.NewOutboxEventsQ(db),
			profiles: pgdb.NewProfilesQ(db),
		},
	}
}
