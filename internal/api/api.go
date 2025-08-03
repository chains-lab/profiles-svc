package api

import (
	"context"
	"fmt"
	"net"

	"github.com/chains-lab/profiles-proto/gen/go/svc"
	"github.com/chains-lab/profiles-svc/internal/api/interceptors"
	"github.com/chains-lab/profiles-svc/internal/api/service"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/config"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg config.Config, log logger.Logger, app *app.App) error {
	server := service.NewService(cfg, app)
	authInterceptor := interceptors.NewAuth(cfg.JWT.Service.SecretKey, cfg.JWT.User.AccessToken.SecretKey)
	logInterceptor := logger.UnaryLogInterceptor(log)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logInterceptor,
			authInterceptor,
		),
	)

	svc.RegisterProfilesServer(grpcServer, server)

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	log.Infof("gRPC server listening on %s", lis.Addr())

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- grpcServer.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down gRPC server …")
		grpcServer.GracefulStop()
		return nil
	case err := <-serveErrCh:
		return fmt.Errorf("gRPC Serve() exited: %w", err)
	}
}
