package callback

import (
	"context"
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

const (
	InboxStatusPending   = "pending"
	InboxStatusProcessed = "processed"
)

func (s Service) CreateAccount(ctx context.Context, event kafka.Message) error {
	//var p AccountCreatedPayload
	//
	//if len(event.Value) == 0 {
	//	return fmt.Errorf("empty kafka message value")
	//}
	//if err := json.Unmarshal(event.Value, &p); err != nil {
	//	return fmt.Errorf("unmarshal AccountCreatedPayload: %w", err)
	//}

	s.log.Infof("received create account event for account %s", string(event.Key))

	msg := contracts.Message{
		Topic:        contracts.AccountsTopicV1,
		EventType:    AccountCreatedEvent,
		EventVersion: 1,
		Key:          string(event.Key),
		Payload:      event.Value,
	}

	//profile, err := s.domain.CreateProfile(ctx, p.Account.ID, p.Account.Username)
	//if err != nil {
	//	s.log.Errorf("failed to create profile for account %s: %v", p.Account.ID, err)
	//}
	//
	//status := InboxStatusProcessed
	//if profile.IsNil() {
	//	status = InboxStatusPending
	//}

	err := s.inbox.CreateInboxEvent(ctx, contracts.InboxEvent{
		ID:           uuid.New(),
		Topic:        msg.Topic,
		EventType:    msg.EventType,
		EventVersion: msg.EventVersion,
		Key:          msg.Key,
		Payload:      msg.Payload,
		Status:       InboxStatusPending,
		NextRetryAt:  time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
	})
	if err != nil {
		s.log.Errorf("failed to upsert inbox event for account %s: %v", string(event.Key), err)
		return fmt.Errorf("failed to processing create account event for account %s: %w", string(event.Key), err)
	}

	return nil
}
