package models

import (
	"time"

	"github.com/google/uuid"
)

type Residence struct {
	UserID    uuid.UUID  `json:"user_id"`
	Country   *string    `json:"country,omitempty"`
	City      *string    `json:"city,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
