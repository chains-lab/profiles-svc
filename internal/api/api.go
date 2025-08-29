package api

import (
	"context"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/api/rest"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/config"
)

func Start(ctx context.Context, cfg config.Config, log logium.Logger, wg *sync.WaitGroup, app *app.App) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	restApi := rest.NewRest(cfg, log, app)
	run(func() { restApi.Run(ctx) })
}
