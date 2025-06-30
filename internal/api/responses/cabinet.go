package responses

import (
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/electorcab"
)

func Cabinet(model models.Cabinet) *electorcab.Cabinet {
	bio := Biography(model.Biography)
	job := JobResume(model.Job)
	profile := Profile(model.Profile)

	return &electorcab.Cabinet{
		Profile:   profile,
		Biography: bio,
		JobResume: job,
	}
}
