package responses

import (
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func JobResume(model models.JobResume) *electorcab.JobResume {
	job := electorcab.JobResume{
		UserId:   model.UserID.String(),
		Degree:   model.Degree,
		Industry: model.Industry,
		Income:   model.Income,
	}

	if model.DegreeUpdatedAt != nil {
		upAt := model.DegreeUpdatedAt.String()
		job.DegreeUpdatedAt = &upAt
	}

	if model.IndustryUpdatedAt != nil {
		upAt := model.IndustryUpdatedAt.String()
		job.IndustryUpdatedAt = &upAt
	}

	if model.IncomeUpdatedAt != nil {
		upAt := model.IncomeUpdatedAt.String()
		job.IncomeUpdatedAt = &upAt
	}

	return &job
}
