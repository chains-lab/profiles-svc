package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/profiles-svc/internal/events/contracts"
	"github.com/chains-lab/profiles-svc/internal/repo/pgdb"
	"github.com/google/uuid"
)

func (r *Repository) CreateInboxEvent(
	ctx context.Context,
	event contracts.InboxEvent,
) error {
	err := r.sql.inbox.New().Insert(ctx, pgdb.InboxEvent{
		ID:           event.ID,
		Topic:        event.Topic,
		EventType:    event.EventType,
		EventVersion: event.EventVersion,
		Key:          event.Key,
		Payload:      event.Payload,
		Status:       event.Status,
		NextRetryAt:  event.NextRetryAt,
		CreatedAt:    event.CreatedAt,
		ProcessedAt:  event.ProcessedAt,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetInboxEvent(
	ctx context.Context,
	id uuid.UUID,
) (contracts.InboxEvent, error) {
	res, err := r.sql.inbox.New().FilterID(id).Get(ctx)
	switch {
	case err != nil:
		return contracts.InboxEvent{}, err
	case errors.Is(err, sql.ErrNoRows):
		return contracts.InboxEvent{}, nil
	}

	e := res.ToEntity()
	return e, nil
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
	return r.sql.inbox.FilterID(ids...).AddAttempts().UpdateStatus("processed").Update(ctx)
}

func (r *Repository) DelayInboxEvents(
	ctx context.Context,
	ids []uuid.UUID,
	delay time.Duration,
) error {
	nextRetryAt := time.Now().UTC().Add(delay)
	return r.sql.inbox.FilterID(ids...).AddAttempts().UpdateNextRetryAt(nextRetryAt).Update(ctx)
}
