package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/profiles-svc/internal/api/rest/meta"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	GetOwnProfile(w http.ResponseWriter, r *http.Request)
	CreateOwnProfile(w http.ResponseWriter, r *http.Request)
	UpdateOwnProfile(w http.ResponseWriter, r *http.Request)
	UpdateOwnUsername(w http.ResponseWriter, r *http.Request)

	GetProfileByUsername(w http.ResponseWriter, r *http.Request)
	GetProfileByID(w http.ResponseWriter, r *http.Request)

	ResetUsername(w http.ResponseWriter, r *http.Request)
	ResetProfile(w http.ResponseWriter, r *http.Request)

	UpdateOfficial(w http.ResponseWriter, r *http.Request)
}

func (s *Service) Run(ctx context.Context, h Handlers) {
	svcAuth := mdlv.ServiceGrant(enum.SsoSVC, s.cfg.JWT.Service.SecretKey)
	userAuth := mdlv.Auth(meta.UserCtxKey, s.cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin:     true,
		roles.SuperUser: true,
	})

	s.router.Route("/profiles-svc", func(r chi.Router) {
		r.Use(svcAuth)
		r.Route("/v1", func(r chi.Router) {
			r.Route("/profiles", func(r chi.Router) {
				r.With(userAuth).Route("/own", func(r chi.Router) {
					r.Get("/", h.GetOwnProfile)
					r.Post("/", h.CreateOwnProfile)
					r.Patch("/", h.UpdateOwnProfile)
					r.Patch("/username", h.UpdateOwnUsername)
				})

				r.Get("/username/{username}", h.GetProfileByUsername)
				r.Get("/user_id/{user_id}", h.GetProfileByID)

				r.With(sysadmin).Route("/admin", func(r chi.Router) {
					r.Route("/{user_id}", func(r chi.Router) {
						r.Route("/reset", func(r chi.Router) {
							r.Post("/username", h.ResetUsername)
							r.Post("/profile", h.ResetProfile)
						})

						r.Patch("/official/{value}", h.UpdateOfficial)
					})
				})
			})
		})
	})

	s.Start(ctx)

	<-ctx.Done()
	s.Stop(ctx)
}
