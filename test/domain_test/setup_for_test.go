package domain_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/chains-lab/profiles-svc/internal"
	"github.com/chains-lab/profiles-svc/internal/data"
	"github.com/chains-lab/profiles-svc/internal/domain/services/profile"
)

// TEST DATABASE CONNECTION
const testDatabaseURL = "postgresql://postgres:postgres@localhost:7777/postgres?sslmode=disable"

func mustExec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatalf("exec failed: %v", err)
	}
}

type Setup struct {
	domain domain
	Cfg    internal.Config
}

type domain struct {
	profile profile.Service
}

func newSetup(t *testing.T) (Setup, error) {
	cfg := internal.Config{
		Database: internal.DatabaseConfig{
			SQL: struct {
				URL string `mapstructure:"url"`
			}{
				URL: testDatabaseURL,
			},
		},
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := data.NewDatabase(pg)

	profileSvc := profile.New(database)

	return Setup{
		domain: domain{
			profile: profileSvc,
		},
		Cfg: cfg,
	}, nil
}
