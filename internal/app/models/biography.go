package models

import (
	"time"

	"github.com/google/uuid"
)

type Biography struct {
	UserID   uuid.UUID  `json:"userid"`
	Birthday *time.Time `json:"birthday"`
	Sex      *string    `json:"sex"`
	City     *string    `json:"city,omitempty"`
	Region   *string    `json:"region,omitempty"`
	Country  *string    `json:"country,omitempty"`

	SexUpdatedAt       *time.Time `json:"sex_updated_at"`
	ResidenceUpdatedAt *time.Time `json:"residence_updated_at,omitempty"`
}
