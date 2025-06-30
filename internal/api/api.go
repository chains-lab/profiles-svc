package api

import (
	"context"
	"fmt"
	"net"

	"github.com/chains-lab/elector-cab-svc/internal/api/interceptors"
	"github.com/chains-lab/elector-cab-svc/internal/api/service"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	"github.com/chains-lab/elector-cab-svc/internal/config"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService interface {
	CreateCabinet(context.Context, *emptypb.Empty) (*svc.Cabinet, error)

	GetOwnCabinet(context.Context, *emptypb.Empty) (*svc.Cabinet, error)
	GetOwnProfile(context.Context, *emptypb.Empty) (*svc.Profile, error)
	GetOwnBiography(context.Context, *emptypb.Empty) (*svc.Biography, error)
	GetOwnJobResume(context.Context, *emptypb.Empty) (*svc.JobResume, error)
	// Profile
	UpdateOwnProfile(context.Context, *svc.UpdateOwnProfileRequest) (*svc.Profile, error)
	// Biography
	UpdateOwnSex(context.Context, *svc.UpdateOwnSexRequest) (*svc.Biography, error)
	UpdateOwnBirthday(context.Context, *svc.UpdateOwnBirthdayRequest) (*svc.Biography, error)
	UpdateOwnNationality(context.Context, *svc.UpdateOwnNationalityRequest) (*svc.Biography, error)
	UpdateOwnPrimaryLanguage(context.Context, *svc.UpdateOwnPrimaryLanguageRequest) (*svc.Biography, error)
	UpdateOwnResidence(context.Context, *svc.UpdateOwnResidenceRequest) (*svc.Biography, error)
	// Job
	UpdateOwnDegree(context.Context, *svc.UpdateOwnDegreeRequest) (*svc.JobResume, error)
	UpdateOwnIndustry(context.Context, *svc.UpdateOwnIndustryRequest) (*svc.JobResume, error)
	UpdateOwnIncome(context.Context, *svc.UpdateOwnIncomeRequest) (*svc.JobResume, error)
}

type AdminService interface {
	AdminGetCabinet(context.Context, *svc.AdminGetCabinetRequest) (*svc.Cabinet, error)
	AdminGetProfile(context.Context, *svc.AdminGetProfileRequest) (*svc.Profile, error)
	AdminUpdateProfile(context.Context, *svc.AdminUpdateProfileRequest) (*svc.Profile, error)
	AdminGetBiography(context.Context, *svc.AdminGetBiographyRequest) (*svc.Biography, error)
	AdminUpdateBiography(context.Context, *svc.AdminUpdateBiographyRequest) (*svc.Biography, error)
	AdminGetJobResume(context.Context, *svc.AdminGetJobResumeRequest) (*svc.JobResume, error)
	AdminUpdateJobResume(context.Context, *svc.AdminUpdateJobResumeRequest) (*svc.JobResume, error)
}

func Run(ctx context.Context, cfg config.Config, log *logrus.Logger, app *app.App) error {
	// 1) Создаём реализацию хэндлеров и interceptor
	server := service.NewService(cfg, app)
	authInterceptor := interceptors.NewAuth(cfg.JWT.Service.SecretKey)

	// 2) Инициализируем gRPC‐сервер
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)
	svc.RegisterUserServiceServer(grpcServer, server)
	svc.RegisterAdminServiceServer(grpcServer, server)

	// 3) Открываем слушатель
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	log.Infof("gRPC server listening on %s", lis.Addr())

	// 4) Запускаем Serve в горутине
	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- grpcServer.Serve(lis)
	}()

	// 5) Слушаем контекст и окончание Serve()
	select {
	case <-ctx.Done():
		log.Info("shutting down gRPC server …")
		grpcServer.GracefulStop()
		return nil
	case err := <-serveErrCh:
		return fmt.Errorf("gRPC Serve() exited: %w", err)
	}
}
