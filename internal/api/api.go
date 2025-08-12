package api

import (
	"context"
	"fmt"
	"net"

	profilesProto "github.com/chains-lab/profiles-proto/gen/go/svc/profile"
	profilesAdminProto "github.com/chains-lab/profiles-proto/gen/go/svc/profileadmin"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/interceptor"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/service/profile"
	"github.com/chains-lab/profiles-svc/internal/api/grpc/service/profileadmin"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/config"
	"github.com/chains-lab/profiles-svc/internal/logger"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg config.Config, log logger.Logger, app *app.App) error {
	profileSVC := profile.NewService(cfg, app)
	profileAdminSVC := profileadmin.NewService(cfg, app)
	authInterceptor := interceptor.Auth(cfg.JWT.Service.SecretKey)
	logInterceptor := logger.UnaryLogInterceptor(log)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logInterceptor,
			authInterceptor,
		),
	)

	profilesProto.RegisterProfilesServiceServer(grpcServer, profileSVC)
	profilesAdminProto.RegisterProfileAdminServiceServer(grpcServer, profileAdminSVC)

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
