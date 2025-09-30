package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/chains-lab/profiles-svc/internal/rest/requests"
	"github.com/chains-lab/profiles-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) UpdateOfficial(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid user id")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	req, err := requests.UpdateOfficial(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid update official request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.Profile.UpdateOfficial(r.Context(), userID, req.Data.Attributes.Official)
	if err != nil {
		s.log.WithError(err).Errorf("failed to update official status")
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
