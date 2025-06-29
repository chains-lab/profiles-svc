package models

import (
	"time"

	"github.com/google/uuid"
)

type JobResume struct {
	UserID   uuid.UUID `json:"user_id"`
	Degree   *string   `json:"degree"`
	Industry *string   `json:"industry"`
	Income   *string   `json:"income"`

	DegreeUpdatedAt   *time.Time `json:"degree_updated_at,omitempty"`
	IndustryUpdatedAt *time.Time `json:"industry_updated_at,omitempty"`
	IncomeUpdatedAt   *time.Time `json:"income_updated_at,omitempty"`
}
