package responses

import (
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func Biography(model models.Biography) *electorcab.Biography {
	bio := electorcab.Biography{
		UserId:  model.UserID.String(),
		Sex:     model.Sex,
		Country: model.Country,
		City:    model.City,
	}

	if model.Birthday != nil {
		hb := model.Birthday.String()
		bio.Birthday = &hb
	}

	if model.SexUpdatedAt != nil {
		upAt := model.SexUpdatedAt.String()
		bio.SexUpdatedAt = &upAt
	}

	if model.ResidenceUpdatedAt != nil {
		upAt := model.ResidenceUpdatedAt.String()
		bio.ResidenceUpdatedAt = &upAt
	}

	return &bio
}
