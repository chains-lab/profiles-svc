package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/profiles-svc/internal/api/rest/meta"
	"github.com/chains-lab/profiles-svc/internal/api/rest/requests"
	"github.com/chains-lab/profiles-svc/internal/api/rest/responses"
	"github.com/chains-lab/profiles-svc/internal/errx"
)

func (s Service) UpdateOwnUsername(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateUsername(r)
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid update username request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.app.UpdateUsername(r.Context(), initiator.UserID, req.Data.Attributes.Username)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to update username")
		switch {
		case errors.Is(err, errx.ErrorProfileForUserDoesNotExist):
			ape.RenderErr(w, problems.NotFound("profile for user does not exist"))
		case errors.Is(err, errx.ErrorUserTooYoung):
			ape.RenderErr(w, problems.Unauthorized("birthday must be at least 12 years ago"))
		case errors.Is(err, errx.ErrorSexIsNotValid):
			ape.RenderErr(w, problems.InvalidPointer("/data/attributes/sex", err))
		case errors.Is(err, errx.ErrorBirthdateIsNotValid):
			ape.RenderErr(w, problems.InvalidPointer("/data/attributes/birthdate", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Profile(res))
}
