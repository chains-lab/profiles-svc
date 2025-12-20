package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/umisto/kafkakit/box"
	"github.com/umisto/logium"
	"github.com/umisto/profiles-svc/internal"
	"github.com/umisto/profiles-svc/internal/domain/modules/profile"
	"github.com/umisto/profiles-svc/internal/events/consumer"
	"github.com/umisto/profiles-svc/internal/events/consumer/callback"
	"github.com/umisto/profiles-svc/internal/repo"
	"github.com/umisto/profiles-svc/internal/rest/middlewares"

	"github.com/umisto/profiles-svc/internal/rest"
	"github.com/umisto/profiles-svc/internal/rest/controller"
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

	database := repo.New(pg)
	kafkaBox := box.New(pg)

	profileSvc := profile.New(database)

	ctrl := controller.New(log, profileSvc)
	mdlv := middlewares.New(log)

	kafkaConsumer := consumer.New(log, cfg.Kafka.Brokers, callback.NewService(log, kafkaBox))
	kafkaInboxWorker := consumer.NewInboxWorker(log, kafkaBox, profileSvc)

	run(func() { kafkaConsumer.Run(ctx) })

	run(func() { kafkaInboxWorker.Run(ctx) })

	run(func() { rest.Run(ctx, cfg, log, mdlv, ctrl) })
}
