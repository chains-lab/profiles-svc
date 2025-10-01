package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal"
	"github.com/chains-lab/profiles-svc/internal/rest/meta"
	"github.com/go-chi/chi/v5"
)

type Controllers interface {
	GetOwnProfile(w http.ResponseWriter, r *http.Request)

	CreateOwnProfile(w http.ResponseWriter, r *http.Request)
	GetProfileByUsername(w http.ResponseWriter, r *http.Request)
	GetProfileByID(w http.ResponseWriter, r *http.Request)

	FilterProfiles(w http.ResponseWriter, r *http.Request)

	UpdateOwnProfile(w http.ResponseWriter, r *http.Request)
	UpdateOwnUsername(w http.ResponseWriter, r *http.Request)
	UpdateOwnSex(w http.ResponseWriter, r *http.Request)
	UpdateOwnBirthDate(w http.ResponseWriter, r *http.Request)
	UpdateOfficial(w http.ResponseWriter, r *http.Request)

	ResetUsername(w http.ResponseWriter, r *http.Request)
	ResetProfile(w http.ResponseWriter, r *http.Request)
}

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, c Controllers) {
	svcAuth := mdlv.ServiceGrant(enum.SsoSVC, cfg.JWT.Service.SecretKey)
	userAuth := mdlv.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Moder: true,
		roles.Admin: true,
	})

	r := chi.NewRouter()

	r.Route("/profile-svc", func(r chi.Router) {
		r.Use(svcAuth)
		r.Route("/v1", func(r chi.Router) {
			r.Route("/profile", func(r chi.Router) {
				r.Get("/", c.FilterProfiles)

				r.Get("/username/{username}", c.GetProfileByUsername)
				r.Get("/user_id/{user_id}", c.GetProfileByID)

				r.With(userAuth).Route("/own", func(r chi.Router) {
					r.Get("/", c.GetOwnProfile)
					r.Post("/", c.CreateOwnProfile)

					r.Put("/", c.UpdateOwnProfile)
					r.Patch("/sex", c.UpdateOwnSex)
					r.Patch("/username", c.UpdateOwnUsername)
					r.Patch("/birth_date", c.UpdateOwnBirthDate)
				})

				r.With(sysadmin).Route("/admin", func(r chi.Router) {
					r.Route("/{user_id}", func(r chi.Router) {
						r.Route("/reset", func(r chi.Router) {
							r.Post("/username", c.ResetUsername)
							r.Post("/profile", c.ResetProfile)
						})

						r.Patch("/official/{value}", c.UpdateOfficial)
					})
				})
			})
		})
	})

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	<-ctx.Done()

	log.Info("shutting down REST service")
}
