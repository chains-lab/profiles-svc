package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/profiles-svc/internal/api/rest/responses"
	"github.com/chains-lab/profiles-svc/internal/errx"
	"github.com/go-chi/chi/v5"
)

func (s Service) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	res, err := s.app.GetProfileByUsername(r.Context(), username)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to get profile by username")
		switch {
		case errors.Is(err, errx.ErrorProfileForUserDoesNotExist):
			ape.RenderErr(w, problems.NotFound("profile for user does not exist"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Profile(res))
}
