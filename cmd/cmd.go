package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal"
	"github.com/chains-lab/profiles-svc/internal/data"
	"github.com/chains-lab/profiles-svc/internal/domain/services/profile"

	"github.com/chains-lab/profiles-svc/internal/rest"
	"github.com/chains-lab/profiles-svc/internal/rest/controller"
)

func StartServices(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := data.NewDatabase(pg)

	profileSvc := profile.New(database)

	ctrl := controller.New(log, profileSvc)

	run(func() { rest.Run(ctx, cfg, log, ctrl) })
}
