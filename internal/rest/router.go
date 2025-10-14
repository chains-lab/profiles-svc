package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal"
	"github.com/chains-lab/profiles-svc/internal/rest/meta"
	"github.com/chains-lab/restkit/roles"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	GetMyProfile(w http.ResponseWriter, r *http.Request)

	CreateMyProfile(w http.ResponseWriter, r *http.Request)
	GetProfileByUsername(w http.ResponseWriter, r *http.Request)
	GetProfileByID(w http.ResponseWriter, r *http.Request)

	FilterProfiles(w http.ResponseWriter, r *http.Request)

	UpdateMyProfile(w http.ResponseWriter, r *http.Request)
	UpdateMyUsername(w http.ResponseWriter, r *http.Request)
	UpdateOfficial(w http.ResponseWriter, r *http.Request)

	ResetProfile(w http.ResponseWriter, r *http.Request)
}

type Middleware interface {
	ServiceGrant(serviceName, skService string) func(http.Handler) http.Handler
	Auth(userCtxKey interface{}, skUser string) func(http.Handler) http.Handler
	RoleGrant(userCtxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler
}

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, m Middleware, h Handlers) {
	svcAuth := m.ServiceGrant(cfg.Service.Name, cfg.JWT.Service.SecretKey)
	userAuth := m.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysmoder := m.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Moder: true,
		roles.Admin: true,
	})

	r := chi.NewRouter()

	r.Route("/profile-svc", func(r chi.Router) {
		r.Use(svcAuth)
		r.Route("/v1", func(r chi.Router) {
			r.Route("/profile", func(r chi.Router) {
				r.Get("/", h.FilterProfiles)

				r.Get("/u/{username}", h.GetProfileByUsername)

				r.With(userAuth).Route("/me", func(r chi.Router) {
					r.Post("/", h.CreateMyProfile)

					r.Get("/", h.GetMyProfile)
					r.Put("/", h.UpdateMyProfile)

					r.Patch("/username", h.UpdateMyUsername)
				})

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", h.GetProfileByID)

					r.With(sysmoder).Patch("/official", h.UpdateOfficial)
					r.With(sysmoder).Put("/reset", h.ResetProfile)
				})
			})
		})
	})

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	<-ctx.Done()

	log.Info("shutting dMy REST service")
}
