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

func (p Profile) IsNil() bool {
	return p.AccountID == uuid.Nil
}

type ProfileCollection struct {
	Data  []Profile `json:"data"`
	Page  int32     `json:"page"`
	Size  int32     `json:"size"`
	Total int32     `json:"total"`
}
