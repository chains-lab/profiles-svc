package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/chains-lab/profiles-svc/internal/events/consumer/callback"
	"github.com/chains-lab/profiles-svc/internal/events/consumer/subscriber"
	"github.com/chains-lab/profiles-svc/internal/events/contracts"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	addr      string
	callbacks callbacks
	inbox     inbox
	domain    domain
	log       logium.Logger
}

type callbacks interface {
	CreateAccount(ctx context.Context, event kafka.Message) error
	UpdateUsername(ctx context.Context, event kafka.Message) error
}

type inbox interface {
	GetPendingInboxEvents(
		ctx context.Context,
		limit uint,
	) ([]contracts.InboxEvent, error)

	MarkInboxEventsAsProcessed(
		ctx context.Context,
		ids []uuid.UUID,
	) error

	DelayInboxEvents(
		ctx context.Context,
		ids []uuid.UUID,
		delay time.Duration,
	) error
}

type domain interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, username string) (entity.Profile, error)
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (entity.Profile, error)
}

func New(log logium.Logger, addr string, inbox inbox, callbacks callbacks, domain domain) *Service {
	return &Service{
		addr:      addr,
		log:       log,
		callbacks: callbacks,
		inbox:     inbox,
		domain:    domain,
	}
}

func (s Service) Run(ctx context.Context) {
	sub := subscriber.New(s.addr, contracts.AccountsTopicV1, contracts.GroupProfilesSvc)

	s.log.Info("starting events consumer", "addr", s.addr)

	go func() {
		err := sub.Consume(ctx, func(m kafka.Message) (subscriber.HandlerFunc, bool) {
			et, ok := subscriber.Header(m, "event_type")
			if !ok {
				return nil, false
			}

			switch et {
			case "account.created":
				return s.callbacks.CreateAccount, true
			case "account.username.change":
				return s.callbacks.UpdateUsername, true
			default:
				return nil, false
			}
		})
		if err != nil {
			log.Printf("accounts consumer stopped: %v", err)
		}
	}()
}

const eventInboxRetryDelay = 1 * time.Minute

func (s Service) InboxWorker(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		events, err := s.inbox.GetPendingInboxEvents(ctx, 10)
		if err != nil {
			s.log.Error("failed to get pending inbox events", "error", err)
			continue
		}
		if len(events) == 0 {
			continue
		}

		var processed []uuid.UUID
		var delayed []uuid.UUID

		for _, ev := range events {
			s.log.Info("processing inbox event", "id", ev.ID, "type", ev.EventType)

			key, err := uuid.Parse(ev.Key)
			if err != nil {
				s.log.Error("invalid inbox event key", "id", ev.ID, "error", err)
				processed = append(processed, ev.ID)
				continue
			}

			switch ev.EventType {
			case "account.created":
				var p callback.AccountCreatedPayload
				if err := json.Unmarshal(ev.Payload, &p); err != nil {
					s.log.Error("bad payload for account.create", "id", ev.ID, "error", err)
					processed = append(processed, ev.ID)
					continue
				}
				if _, err := s.domain.CreateProfile(ctx, key, p.Account.Username); err != nil {
					s.log.Error("failed to create profile", "id", ev.ID, "error", err)
					delayed = append(delayed, ev.ID)
					continue
				}
				processed = append(processed, ev.ID)

			case "account.username.change":
				var p callback.AccountUsernameChangePayload
				if err := json.Unmarshal(ev.Payload, &p); err != nil {
					s.log.Error("bad payload for account.username.change", "id", ev.ID, "error", err)
					processed = append(processed, ev.ID)
					continue
				}
				if _, err := s.domain.UpdateProfileUsername(ctx, key, p.Account.Username); err != nil {
					s.log.Error("failed to update username", "id", ev.ID, "error", err)
					delayed = append(delayed, ev.ID)
					continue
				}
				processed = append(processed, ev.ID)

			default:
				s.log.Warn("unknown inbox event type", "id", ev.ID, "type", ev.EventType)
				processed = append(processed, ev.ID)
			}
		}

		if len(processed) > 0 {
			if err := s.inbox.MarkInboxEventsAsProcessed(ctx, processed); err != nil {
				s.log.Error("failed to mark inbox events as processed", "error", err)
			}
		}
		if len(delayed) > 0 {
			if err := s.inbox.DelayInboxEvents(ctx, delayed, eventInboxRetryDelay); err != nil {
				s.log.Error("failed to delay inbox events", "error", err)
			}
		}
	}
}
