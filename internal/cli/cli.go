package cli

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/chains-lab/profile-storage/internal/app"
	"github.com/chains-lab/profile-storage/internal/config"
)

func Run(args []string) bool {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := config.SetupLogger(cfg.Server.Log.Level, cfg.Server.Log.Format)
	logger.Info("Starting server...")

	var (
		service = kingpin.New("chains-auth", "")
		runCmd  = service.Command("run", "run command")

		serviceCmd     = runCmd.Command("service", "run service")
		migrateCmd     = service.Command("migrate", "migrate command")
		migrateUpCmd   = migrateCmd.Command("up", "migrate db up")
		migrateDownCmd = migrateCmd.Command("down", "migrate db down")

		//docs = service.Command("docs", "documentation command")
		//
		//generateDocs = docs.Command("generate", "generate API documentation")
		//openDocs     = docs.Command("open", "open API documentation in browser")
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Fatalf("failed to create server: %v", err)
		return false
	}

	var wg sync.WaitGroup

	cmd, err := service.Parse(args[1:])
	if err != nil {
		logger.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch cmd {
	case serviceCmd.FullCommand():
		runServices(ctx, cfg, logger, &wg, &application)
	case migrateUpCmd.FullCommand():
		err = MigrateUp(ctx, cfg)
	case migrateDownCmd.FullCommand():
		err = MigrateDown(ctx, cfg)
	//case generateDocs.FullCommand():
	//	err = GenerateDocs(ctx, cfg)
	//case openDocs.FullCommand():
	//	err = OpenDocs(ctx, cfg)
	default:
		logger.Errorf("unknown command %s", cmd)
		return false
	}
	if err != nil {
		logger.WithError(err).Error("failed to exec cmd")
		return false
	}

	wgch := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgch)
	}()

	select {
	case <-ctx.Done():
		log.Printf("Interrupt signal received: %v", ctx.Err())
		<-wgch
	case <-wgch:
		log.Print("All services stopped")
	}

	return true
}
