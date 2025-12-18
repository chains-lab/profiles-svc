package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chains-lab/kafkakit/box"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/chains-lab/profiles-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type InboxWorker struct {
	log    logium.Logger
	inbox  inbox
	domain domain
}

type inbox interface {
	GetInboxEventByID(
		ctx context.Context,
		id uuid.UUID,
	) (box.InboxEvent, error)

	GetPendingInboxEvents(
		ctx context.Context,
		limit int32,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsProcessed(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsFailed(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsPending(
		ctx context.Context,
		ids []uuid.UUID,
		delay time.Duration,
	) ([]box.InboxEvent, error)
}

type domain interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, username string) (entity.Profile, error)
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (entity.Profile, error)
}

func NewInboxWorker(
	log logium.Logger,
	inbox inbox,
	domain domain,
) InboxWorker {
	return InboxWorker{
		log:    log,
		inbox:  inbox,
		domain: domain,
	}
}

const eventInboxRetryDelay = 1 * time.Minute

func (w InboxWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		events, err := w.inbox.GetPendingInboxEvents(ctx, 10)
		if err != nil {
			w.log.Error("failed to get pending inbox events", "error", err)
			continue
		}
		if len(events) == 0 {
			continue
		}

		var processed []uuid.UUID
		var delayed []uuid.UUID

		for _, ev := range events {
			w.log.Infof("processing inbox event: %s, type %s", ev.ID, ev.Type)

			key, err := uuid.Parse(ev.Key)
			if err != nil {
				w.log.Error("invalid inbox event key", "id", ev.ID, "error", err)
				processed = append(processed, ev.ID)
				continue
			}

			switch ev.Type {
			case contracts.AccountCreatedEvent:
				var p contracts.AccountCreatedPayload
				if err = json.Unmarshal(ev.Payload, &p); err != nil {
					w.log.Error("bad payload for account.create", "id", ev.ID, "error", err)
					processed = append(processed, ev.ID)
					continue
				}

				if _, err = w.domain.CreateProfile(ctx, key, p.Account.Username); err != nil {
					w.log.Error("failed to create profile", "id", ev.ID, "error", err)
					delayed = append(delayed, ev.ID)
					continue
				}
				processed = append(processed, ev.ID)

			case contracts.AccountUsernameChangeEvent:
				var p contracts.AccountUsernameChangePayload
				if err = json.Unmarshal(ev.Payload, &p); err != nil {
					w.log.Error("bad payload for account.username.change", "id", ev.ID, "error", err)
					processed = append(processed, ev.ID)
					continue
				}

				if _, err = w.domain.UpdateProfileUsername(ctx, key, p.Account.Username); err != nil {
					w.log.Error("failed to update username", "id", ev.ID, "error", err)
					delayed = append(delayed, ev.ID)
					continue
				}
				processed = append(processed, ev.ID)

			default:
				w.log.Warn("unknown inbox event type", "id", ev.ID, "type", ev.Type)
				processed = append(processed, ev.ID)
			}
		}

		if len(processed) > 0 {
			_, err = w.inbox.MarkInboxEventsAsProcessed(ctx, processed)
			if err != nil {
				w.log.Error("failed to mark inbox events as processed", "error", err)
			}
		}

		if len(delayed) > 0 {
			_, err = w.inbox.MarkInboxEventsAsPending(ctx, delayed, eventInboxRetryDelay)
			if err != nil {
				w.log.Error("failed to delay inbox events", "error", err)
			}
		}
	}
}
