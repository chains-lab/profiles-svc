package responses

import (
	"github.com/chains-lab/profiles-proto/gen/go/svc"
	"github.com/chains-lab/profiles-svc/internal/app/models"
)

func Profile(model models.Profile) *svc.Profile {
	var birthdate string
	if model.BirthDate != nil {
		birthdate = model.BirthDate.String()
	}

	return &svc.Profile{
		UserId:      model.UserID.String(),
		Username:    model.Username,
		Pseudonym:   model.Pseudonym,
		Description: model.Description,
		Avatar:      model.Avatar,
		Official:    model.Official,
		Sex:         model.Sex,
		BirthDate:   &birthdate,
		UpdatedAt:   model.UpdatedAt.String(),
		CreatedAt:   model.CreatedAt.String(),
	}
}
