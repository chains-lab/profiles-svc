package models

import (
	"time"

	"github.com/google/uuid"
)

var SexValues = []string{
	"male", "female", "other",
}

func ValidateSex(s string) bool {
	for _, v := range SexValues {
		if s == v {
			return true
		}
	}
	return false
}

type Bio struct {
	UserID          uuid.UUID  `json:"userid"`
	Sex             *string    `json:"sex"`
	Birthday        *time.Time `json:"birthday"`
	Citizenship     *string    `json:"citizenship"`
	Nationality     *string    `json:"nationality"`
	PrimaryLanguage *string    `json:"primary_language"`

	SexUpdatedAt             *time.Time `json:"sex_updated_at"`
	CitizenshipUpdatedAt     *time.Time `json:"citizenship_updated_at"`
	NationalityUpdatedAt     *time.Time `json:"nationality_updated_at"`
	PrimaryLanguageUpdatedAt *time.Time `json:"primary_language_updated_at"`
}
