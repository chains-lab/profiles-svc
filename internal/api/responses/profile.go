package responses

import (
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/electorcab"
)

func Profile(model models.Profile) *electorcab.Profile {
	return &electorcab.Profile{
		UserId:      model.UserID.String(),
		Username:    model.Username,
		Pseudonym:   model.Pseudonym,
		Description: model.Description,
		Avatar:      model.Avatar,
		Official:    model.Official,
		UpdatedAt:   model.UpdatedAt.String(),
		CreatedAt:   model.CreatedAt.String(),
	}
}
