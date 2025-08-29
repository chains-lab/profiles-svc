package responses

import (
	"github.com/chains-lab/profiles-svc/internal/app/models"
	"github.com/chains-lab/profiles-svc/resources"
)

func Profile(m models.Profile) resources.Profile {
	resp := resources.Profile{
		Data: resources.ProfileData{
			Id:   m.UserID.String(),
			Type: resources.ProfileType,
			Attributes: resources.ProfileAttributes{
				Username:  m.Username,
				Official:  m.Official,
				UpdatedAt: m.UpdatedAt,
				CreatedAt: m.CreatedAt,
			},
		},
	}

	if m.Pseudonym != nil {
		resp.Data.Attributes.Pseudonym = m.Pseudonym
	}
	if m.Description != nil {
		resp.Data.Attributes.Description = m.Description
	}
	if m.Avatar != nil {
		resp.Data.Attributes.Avatar = m.Avatar
	}
	if m.Sex != nil {
		resp.Data.Attributes.Sex = m.Sex
	}
	if m.Birthdate != nil {
		resp.Data.Attributes.Birthdate = m.Birthdate
	}

	return resp
}
