package models

import (
	"time"

	"github.com/google/uuid"
)

type Biography struct {
	UserID          uuid.UUID  `json:"userid"`
	Sex             *string    `json:"sex"`
	Birthday        *time.Time `json:"birthday"`
	Nationality     *string    `json:"nationality"`
	PrimaryLanguage *string    `json:"primary_language"`
	Country         *string    `json:"country,omitempty"`
	City            *string    `json:"city,omitempty"`

	SexUpdatedAt             *time.Time `json:"sex_updated_at"`
	NationalityUpdatedAt     *time.Time `json:"nationality_updated_at"`
	PrimaryLanguageUpdatedAt *time.Time `json:"primary_language_updated_at"`
	ResidenceUpdatedAt       *time.Time `json:"residence_updated_at,omitempty"`
}
