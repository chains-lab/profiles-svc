package controller

import (
	"net/http"
)

func (s Service) ResetProfile(w http.ResponseWriter, r *http.Request) {
	//userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	//if err != nil {
	//	s.log.WithError(err).Errorf("invalid user id")
	//	ape.RenderErr(w, problems.BadRequest(validation.Errors{
	//		"query": fmt.Errorf("invalid user id: %s", chi.URLParam(r, "user_id")),
	//	})...)
	//
	//	return
	//}
	//
	//req, err := s.domain.Profile.ResetProfile(r.Context(), userID)
	//if err != nil {
	//	s.log.WithError(err).Errorf("failed to reset username")
	//	switch {
	//	case errors.Is(err, errx.ErrorProfileNotFound):
	//		ape.RenderErr(w, problems.NotFound("profile for user does not exist"))
	//	default:
	//		ape.RenderErr(w, problems.InternalError())
	//	}
	//
	//	return
	//}
	//
	//ape.Render(w, http.StatusOK, responses.Profile(req))
}
