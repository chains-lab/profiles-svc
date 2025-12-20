package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
	"github.com/umisto/profiles-svc/internal/domain/errx"
	"github.com/umisto/profiles-svc/internal/rest/responses"
)

func (s Service) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	res, err := s.domain.GetProfileByUsername(r.Context(), username)
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
