package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chains-lab/profiles-svc/internal/events/contracts"
	"github.com/chains-lab/profiles-svc/internal/repo/pgdb"
	"github.com/google/uuid"
)

type CreateInboxEventParams struct {
	Topic        string
	EventType    string
	EventVersion uint
	Key          string
	Payload      json.RawMessage
}

func (r *Repository) CreateInboxEvent(
	ctx context.Context,
	event contracts.Message,
) error {
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	err = r.sql.inbox.New().Insert(ctx, pgdb.InboxEvent{
		ID:           uuid.New(),
		Topic:        event.Topic,
		EventType:    event.EventType,
		EventVersion: event.EventVersion,
		Key:          event.Key,
		Payload:      payloadBytes,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetPendingInboxEvents(
	ctx context.Context,
	limit uint,
) ([]contracts.InboxEvent, error) {
	res, err := r.sql.inbox.New().FilterStatus("pending").Page(limit, 0).Select(ctx)
	if err != nil {
		return nil, err
	}

	events := make([]contracts.InboxEvent, len(res))
	for i, e := range res {
		events[i] = e.ToEntity()
	}
	return events, nil
}

func (r *Repository) MarkInboxEventsAsProcessed(
	ctx context.Context,
	ids []uuid.UUID,
) error {
	return r.sql.inbox.FilterID(ids...).UpdateStatus("processed").Update(ctx)
}

func (r *Repository) DelayInboxEvents(
	ctx context.Context,
	ids []uuid.UUID,
	delay time.Duration,
) error {
	nextRetryAt := time.Now().UTC().Add(delay)
	return r.sql.inbox.FilterID(ids...).UpdateNextRetryAt(nextRetryAt).Update(ctx)
}
