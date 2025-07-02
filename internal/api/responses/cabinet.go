package responses

import (
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func Cabinet(model models.Cabinet) *electorcab.Cabinet {
	bio := Biography(model.Biography)
	job := JobResume(model.Job)

	return &electorcab.Cabinet{
		UserId:    model.UserID.String(),
		Biography: bio,
		JobResume: job,
	}
}
