package callback

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chains-lab/profiles-svc/internal/events/contracts"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type AccountCreatedPayload struct {
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

const AccountCreatedEvent = "account.created"

func (s Service) CreateAccount(ctx context.Context, event kafka.Message) error {
	var p AccountCreatedPayload

	if len(event.Value) == 0 {
		return fmt.Errorf("empty kafka message value")
	}
	if err := json.Unmarshal(event.Value, &p); err != nil {
		return fmt.Errorf("unmarshal AccountCreatedPayload: %w", err)
	}

	err := s.inbox.CreateInboxEvent(ctx, contracts.Message{
		Topic:        event.Topic,
		EventType:    AccountCreatedEvent,
		EventVersion: 1,
		Key:          p.Account.ID.String(),
		Payload:      p,
	})
	if err != nil {
		return err
	}

	return nil
}
