package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/chains-lab/profiles-svc/internal/rest/meta"
	"github.com/chains-lab/profiles-svc/internal/rest/responses"
)

func (s Service) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	res, err := s.domain.Profile.GetByID(r.Context(), initiator.ID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get profile by user id")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			ape.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Profile(res))
}
