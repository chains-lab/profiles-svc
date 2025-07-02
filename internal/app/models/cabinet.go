package models

import "github.com/google/uuid"

type Cabinet struct {
	UserID    uuid.UUID `json:"user_id"`
	Job       JobResume `json:"job"`
	Biography Biography `json:"biography"`
}
