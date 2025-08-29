package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	UserID      uuid.UUID  `json:"userid"`
	Username    string     `json:"username"`
	Official    bool       `json:"official"`
	Pseudonym   *string    `json:"pseudonym,omitempty"`
	Description *string    `json:"description,omitempty"`
	Avatar      *string    `json:"avatar,omitempty"`
	Sex         *string    `json:"sex"`
	Birthdate   *time.Time `json:"birthdate,omitempty"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
