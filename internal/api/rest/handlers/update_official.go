package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/profiles-svc/internal/api/rest/responses"
	"github.com/chains-lab/profiles-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) UpdateOfficial(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid user id")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	value := chi.URLParam(r, "official")
	official := false
	if value == "true" {
		official = true
	}

	res, err := s.app.UpdateProfileOfficial(r.Context(), userID, official)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to update official status")
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
