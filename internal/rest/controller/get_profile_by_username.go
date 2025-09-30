package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/chains-lab/profiles-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
)

func (s Service) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	res, err := s.domain.Profile.GetByUsername(r.Context(), username)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get profile by username")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			ape.RenderErr(w, problems.NotFound("profile for user does not exist"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Profile(res))
}
