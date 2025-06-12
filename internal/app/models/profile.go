package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Pseudonym   *string   `json:"pseudonym,omitempty"`
	Description *string   `json:"description,omitempty"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	Official    bool      `json:"official"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}
