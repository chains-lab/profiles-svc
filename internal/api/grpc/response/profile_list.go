package response

import (
	pagProto "github.com/chains-lab/profiles-proto/gen/go/common/pagination"
	profilesProto "github.com/chains-lab/profiles-proto/gen/go/svc/profile"
	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/chains-lab/profiles-svc/internal/pagination"
)

func ProfileList(models []models.Profile, response pagination.Response) *profilesProto.ProfilesList {
	list := make([]*profilesProto.Profile, len(models))
	for i, model := range models {
		list[i] = &profilesProto.Profile{
			UserId:      model.UserID.String(),
			Username:    model.Username,
			Pseudonym:   model.Pseudonym,
			Description: model.Description,
			Avatar:      model.Avatar,
			Official:    model.Official,
		}
	}

	return &profilesProto.ProfilesList{
		Profiles: list,
		Pagination: &pagProto.Response{
			Page: response.Page,
			Size: response.Size,
		},
	}
}
