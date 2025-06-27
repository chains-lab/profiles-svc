package cli

import (
	"context"
	"sync"

	"github.com/chains-lab/elector-cab-svc/internal/app"
	"github.com/chains-lab/elector-cab-svc/internal/rest"
	"github.com/chains-lab/elector-cab-svc/internal/utils/config"
	"github.com/sirupsen/logrus"
)

func runServices(ctx context.Context, cfg config.Config, log *logrus.Logger, wg *sync.WaitGroup, app *app.App) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	restApi := rest.NewRest(cfg, log, app)
	run(func() { restApi.Run(ctx, log) })
}
