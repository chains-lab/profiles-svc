package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/chains-lab/profiles-svc/internal/events/contracts"
	"github.com/chains-lab/profiles-svc/internal/repo/pgdb"
	"github.com/google/uuid"
)

type CreateInboxEventParams struct {
	Topic        string
	EventType    string
	EventVersion int32
	Key          string
	Payload      json.RawMessage
}

func (r *Repository) CreateInboxEvent(
	ctx context.Context,
	event contracts.InboxEvent,
) error {
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	_, err = r.sql.CreateInboxEvent(ctx, pgdb.CreateInboxEventParams{
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
	limit int32,
) ([]contracts.InboxEvent, error) {
	res, err := r.sql.GetPendingInboxEvents(ctx, limit)
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
	return r.sql.MarkInboxEventsProcessed(ctx, pgdb.MarkInboxEventsProcessedParams{
		Column1: ids,
		ProcessedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now().UTC(),
		},
	})
}

func (r *Repository) DelayInboxEvents(
	ctx context.Context,
	ids []uuid.UUID,
	delay time.Duration,
) error {
	nextRetryAt := time.Now().UTC().Add(delay)
	return r.sql.DelayInboxEvents(ctx, pgdb.DelayInboxEventsParams{
		Column1:     ids,
		NextRetryAt: nextRetryAt,
	})
}
