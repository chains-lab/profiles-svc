package callback

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/profiles-svc/internal/events/contracts"
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

const AccountUsernameChangeEvent = "account.username.change"

func (s Service) UpdateUsername(ctx context.Context, event kafka.Message) error {
	//var p AccountCreatedPayload
	//
	//if len(event.Value) == 0 {
	//	return fmt.Errorf("empty kafka message value")
	//}
	//if err := json.Unmarshal(event.Value, &p); err != nil {
	//	return fmt.Errorf("unmarshal AccountCreatedPayload: %w", err)
	//}

	eventID, err := uuid.Parse(string(event.Key))
	if err != nil {
		s.log.Errorf("invalid event key for account %s: %v", string(event.Key), err)
		return nil
	}

	msg := contracts.Message{
		Topic:        contracts.AccountsTopicV1,
		EventType:    AccountUsernameChangeEvent,
		EventVersion: 1,
		Key:          string(event.Key),
		Payload:      event.Value,
	}

	//profile, err := s.domain.UpdateProfileUsername(ctx, p.Account.ID, p.Account.Username)
	//if err != nil {
	//	s.log.Errorf("failed to update username for account %s: %v", p.Account.ID, err)
	//	return nil
	//}
	//
	//status := InboxStatusProcessed
	//if profile.IsNil() {
	//	status = InboxStatusPending
	//}

	if err = s.inbox.CreateInboxEvent(ctx, contracts.InboxEvent{
		ID:           eventID,
		Topic:        msg.Topic,
		EventType:    msg.EventType,
		EventVersion: msg.EventVersion,
		Key:          msg.Key,
		Payload:      msg.Payload,
		Status:       InboxStatusPending,
		NextRetryAt:  time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
	}); err != nil {
		s.log.Infof("failed to processed account username change for account %s", string(event.Key))
		return fmt.Errorf("failed to processing account username change event for account %s: %w", string(event.Key), err)
	}

	return fmt.Errorf("create inbox event for account %s: %w", string(event.Key), err)
}
