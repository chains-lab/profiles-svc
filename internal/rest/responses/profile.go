package responses

import (
	"github.com/chains-lab/profiles-svc/internal/domain/models"
	"github.com/chains-lab/profiles-svc/resources"
)

func Profile(m models.Profile) resources.Profile {
	resp := resources.Profile{
		Data: resources.ProfileData{
			Id:   m.UserID,
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
	if m.BirthDate != nil {
		resp.Data.Attributes.BirthDate = m.BirthDate
	}

	return resp
}

func ProfileCollection(m models.ProfileCollection) resources.ProfilesCollection {
	resp := resources.ProfilesCollection{
		Data: make([]resources.ProfileData, 0, len(m.Data)),
	}

	for _, el := range m.Data {
		p := Profile(el).Data

		resp.Data = append(resp.Data, p)
	}

	return resp
}
