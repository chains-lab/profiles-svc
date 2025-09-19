package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/profiles-svc/internal/api/rest/meta"
	"github.com/chains-lab/profiles-svc/internal/api/rest/responses"
)

func (s Handler) GetOwnProfile(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	res, err := s.app.GetProfileByUserID(r.Context(), initiator.UserID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get profile by user id")
		switch {
		case err.Error() == "PROFILE_FOR_USER_DOES_NOT_EXIST":
			ape.RenderErr(w, problems.NotFound("profile for user does not exist"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Profile(res))
}
