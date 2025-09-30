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
	BirthDate   *time.Time `json:"birth_date,omitempty"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ProfileCollection struct {
	Data  []Profile `json:"data"`
	Page  uint64    `json:"page"`
	Size  uint64    `json:"size"`
	Total uint64    `json:"total"`
}
