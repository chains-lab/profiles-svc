package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/chains-lab/profiles-svc/internal/domain/services/profile"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/chains-lab/profiles-svc/internal/rest/meta"
	"github.com/chains-lab/profiles-svc/internal/rest/requests"
	"github.com/chains-lab/profiles-svc/internal/rest/responses"
)

func (s Service) UpdateOwnProfile(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateProfile(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid create profile request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Id != initiator.ID {
		s.log.WithError(err).Errorf("id in body and initiastor id mismacht fir update own profile request")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf(
				"id in body: %s and initiastor id: %s mismacht fir update own profile request",
				req.Data.Id,
				initiator.ID,
			),
		})...)
	}

	res, err := s.domain.Profile.Update(r.Context(), initiator.ID, profile.Update{
		Pseudonym:   req.Data.Attributes.Pseudonym,
		Description: req.Data.Attributes.Description,
		Avatar:      req.Data.Attributes.Avatar,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to update profile")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			ape.RenderErr(w, problems.NotFound("profile for user does not exist"))
		case errors.Is(err, errx.ErrorUserTooYoung):
			ape.RenderErr(w, problems.Unauthorized("birthday must be at least 12 years ago"))
		case errors.Is(err, errx.ErrorSexIsNotValid):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"sex": fmt.Errorf("sex value is not supported, %s", err),
			})...)
		case errors.Is(err, errx.ErrorBirthdateIsNotValid):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"birth_date": fmt.Errorf("birth date format is invalid %s", err),
			})...)
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Profile(res))
}
