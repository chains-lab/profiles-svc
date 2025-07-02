package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	UserID      uuid.UUID `json:"userid"`
	Username    string    `json:"username"`
	Pseudonym   *string   `json:"pseudonym,omitempty"`
	Description *string   `json:"description,omitempty"`
	Avatar      *string   `json:"avatar,omitempty"`
	Official    bool      `json:"official"`

	Biography *Biography `json:"biography"`
	JobResume *JobResume `json:"job_resume"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
