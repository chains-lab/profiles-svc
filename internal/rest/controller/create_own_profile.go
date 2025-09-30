package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/chains-lab/profiles-svc/internal/rest/meta"
	"github.com/chains-lab/profiles-svc/internal/rest/requests"
	"github.com/chains-lab/profiles-svc/internal/rest/responses"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s Service) CreateOwnProfile(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateProfile(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid create profile request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.Profile.Create(r.Context(), initiator.ID, req.Data.Attributes.Username)
	if err != nil {
		s.log.WithError(err).Errorf("failed to create profile")
		switch {
		case errors.Is(err, errx.ErrorProfileAlreadyExists):
			ape.RenderErr(w, problems.Conflict("profile for user already exists"))
		case errors.Is(err, errx.ErrorUsernameAlreadyTaken):
			ape.RenderErr(w, problems.Conflict("username already taken"))
		case errors.Is(err, errx.ErrorUsernameIsNotValid):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"username": fmt.Errorf("username %s is not valid", req.Data.Attributes.Username),
			})...)
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusCreated, responses.Profile(res))
}
