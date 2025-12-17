package entity

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	AccountID   uuid.UUID `json:"account_id"`
	Username    string    `json:"username"`
	Official    bool      `json:"official"`
	Pseudonym   *string   `json:"pseudonym,omitempty"`
	Description *string   `json:"description,omitempty"`
	Avatar      *string   `json:"avatar,omitempty"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (e Profile) IsNil() bool {
	return e.AccountID == uuid.Nil
}

type ProfileCollection struct {
	Data  []Profile `json:"data"`
	Page  uint      `json:"page"`
	Size  uint      `json:"size"`
	Total uint      `json:"total"`
}
