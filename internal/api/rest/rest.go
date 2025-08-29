package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/api/rest/handlers"
	"github.com/chains-lab/profiles-svc/internal/api/rest/meta"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/config"
	"github.com/go-chi/chi/v5"
)

type Rest struct {
	server   *http.Server
	router   *chi.Mux
	handlers handlers.Service

	log logium.Logger
	cfg config.Config
}

func NewRest(cfg config.Config, log logium.Logger, app *app.App) Rest {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}
	hands := handlers.NewService(cfg, log, app)

	router.Use()

	return Rest{
		handlers: hands,
		router:   router,
		server:   server,
		log:      log,
		cfg:      cfg,
	}
}

func (a *Rest) Run(ctx context.Context) {
	//svcAuth := mdlv.ServiceAuthMdl(constant.ServiceName, a.cfg.JWT.Service.SecretKey)
	userAuth := mdlv.AuthMdl(meta.UserCtxKey, a.cfg.JWT.User.AccessToken.SecretKey)
	adminGrant := mdlv.AccessGrant(meta.UserCtxKey, roles.Admin, roles.SuperUser)

	a.router.Route("/profiles-svc", func(r chi.Router) {
		//r.Use(svcAuth)
		r.Route("/v1", func(r chi.Router) {
			r.Route("/profiles", func(r chi.Router) {
				r.With(userAuth).Route("/own", func(r chi.Router) {
					r.Get("/", a.handlers.GetOwnProfile)
					r.Post("/", a.handlers.CreateOwnProfile)
					r.Patch("/", a.handlers.UpdateOwnProfile)
					r.Patch("/username", a.handlers.UpdateOwnProfile)
				})

				r.Get("/username/{username}", a.handlers.GetProfileByUsername)
				r.Get("/user_id/{user_id}", a.handlers.GetProfileByID)

				r.With(adminGrant).Route("/admin", func(r chi.Router) {
					r.Route("/{user_id}", func(r chi.Router) {
						r.Route("/reset", func(r chi.Router) {
							r.Post("/username", a.handlers.ResetUsername)
							r.Post("/profile", a.handlers.ResetProfile)
						})

						r.Patch("/official/{value}", a.handlers.UpdateOfficial)
					})
				})
			})
		})
	})

	a.Start(ctx)

	<-ctx.Done()
	a.Stop(ctx)
}

func (a *Rest) Start(ctx context.Context) {
	go func() {
		a.log.Infof("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (a *Rest) Stop(ctx context.Context) {
	a.log.Info("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Errorf("Server shutdown failed: %v", err)
	}
}
