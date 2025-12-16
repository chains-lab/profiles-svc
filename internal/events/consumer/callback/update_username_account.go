package callback

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type AccountUsernameChangePayload struct {
	Account struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Role     string    `json:"role"`
		Status   string    `json:"status"`

		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		UsernameUpdatedAt time.Time `json:"username_name_updated_at"`
	} `json:"account"`
	Email string `json:"email,omitempty"`
}

func (s Service) UpdateUsernameAccount(ctx context.Context, event kafka.Message) error {
	var p AccountCreatedPayload

	if len(event.Value) == 0 {
		return fmt.Errorf("empty kafka message value")
	}
	if err := json.Unmarshal(event.Value, &p); err != nil {
		return fmt.Errorf("unmarshal AccountCreatedPayload: %w", err)
	}

	profile, err := s.domain.UpdateProfileUsername(ctx, p.Account.ID, p.Account.Username)
	if err != nil {
		return err
	}

	s.log.Infof("Profile created successfully: %v", profile)
	return nil
}
