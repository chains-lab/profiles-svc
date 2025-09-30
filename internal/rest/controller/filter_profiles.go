package controller

import (
	"net/http"
	"strings"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/profiles-svc/internal/domain/services/profile"
)

func (s Service) FilterProfiles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	pag, size := pagi.GetPagination(r)

	filters := profile.FilterParams{}

	if username := strings.TrimSpace(q.Get("username")); username != "" {
		filters.Username = &[]string{username}[0]
	}
	if pseudonym := strings.TrimSpace(q.Get("pseudonym")); pseudonym != "" {
		filters.Pseudonym = &[]string{pseudonym}[0]
	}
	if official := strings.TrimSpace(q.Get("official")); official != "" {
		switch official {
		case "true":
			t := true
			filters.Official = &t
		case "false":
			f := false
			filters.Official = &f
		}
	}

	res, err := s.domain.Profile.Filter(ctx)
}
