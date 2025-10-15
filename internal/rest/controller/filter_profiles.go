package controller

import (
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/profiles-svc/internal/domain/services/profile"
	"github.com/chains-lab/profiles-svc/internal/rest/responses"
	"github.com/chains-lab/restkit/pagi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) FilterProfiles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pag, size := pagi.GetPagination(r)

	filters := profile.FilterParams{}

	if userIDs := q["user_id"]; len(userIDs) > 0 {
		for _, raw := range userIDs {
			id, err := uuid.Parse(strings.TrimSpace(raw))
			if err != nil {
				ape.RenderErr(w, problems.BadRequest(validation.Errors{
					"user_id": err,
				})...)
				return
			}
			filters.UserID = append(filters.UserID, id)
		}
	}

	s.log.Debugf("FilterProfiles: parsed user IDs: %v", filters.UserID)

	if usernames := q["username"]; len(usernames) > 0 {
		for _, name := range usernames {
			name = strings.TrimSpace(name)
			if name != "" {
				filters.Username = append(filters.Username, name)
			}
		}
	}

	if usernameLike := strings.TrimSpace(q.Get("username_like")); usernameLike != "" {
		filters.UsernameLike = &usernameLike
	}

	if pseudonym := strings.TrimSpace(q.Get("pseudonym")); pseudonym != "" {
		filters.Pseudonym = &pseudonym
	}

	if official := strings.TrimSpace(q.Get("official")); official != "" {
		switch strings.ToLower(official) {
		case "true":
			t := true
			filters.Official = &t
		case "false":
			f := false
			filters.Official = &f
		}
	}

	res, err := s.domain.Profile.Filter(r.Context(), filters, pag, size)
	if err != nil {
		s.log.WithError(err).Error("failed to filter profiles")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.ProfileCollection(res))
}
