package models

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	UserID   uuid.UUID `json:"user_id"`
	Degree   *string   `json:"degree"`
	Industry *string   `json:"industry"`
	Income   *string   `json:"income"`

	DegreeUpdatedAt   *time.Time `json:"degree_updated_at,omitempty"`
	IndustryUpdatedAt *time.Time `json:"industry_updated_at,omitempty"`
	IncomeUpdatedAt   *time.Time `json:"income_updated_at,omitempty"`
}

var DegreeValues = []string{
	"high_school", "bachelor", "master", "doctorate",
}

func ValidateDegree(d string) bool {
	for _, v := range DegreeValues {
		if d == v {
			return true
		}
	}
	return false
}

var IndustryValues = []string{
	"technology", "healthcare", "finance", "education", "manufacturing", "retail", "other",
}

func ValidateIndustry(i string) bool {
	for _, v := range IndustryValues {
		if i == v {
			return true
		}
	}
	return false
}

var IncomeValues = []string{
	"below_30k", "30k_to_50k", "50k_to_70k", "70k_to_100k", "above_100k",
}

func ValidateIncome(i string) bool {
	for _, v := range IncomeValues {
		if i == v {
			return true
		}
	}
	return false
}
