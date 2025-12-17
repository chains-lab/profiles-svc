package pgdb

import (
	"context"
	"database/sql"

	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/chains-lab/profiles-svc/internal/events/contracts"
)

type txKeyType struct{}

var TxKey = txKeyType{}

func TxFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(TxKey).(*sql.Tx)
	return tx, ok
}

func (p Profile) ToEntity() entity.Profile {
	profile := entity.Profile{
		AccountID:   p.AccountID,
		Username:    p.Username,
		Official:    p.Official,
		Pseudonym:   p.Pseudonym,
		Description: p.Description,
		Avatar:      p.Avatar,

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	return profile
}

func (eo OutboxEvent) ToEntity() contracts.OutboxEvent {
	res := contracts.OutboxEvent{
		ID:           eo.ID,
		Topic:        eo.Topic,
		EventType:    eo.EventType,
		EventVersion: uint(eo.EventVersion),
		Key:          eo.Key,
		Payload:      eo.Payload,
		Status:       eo.Status,
		Attempts:     uint(eo.Attempts),
		NextRetryAt:  eo.NextRetryAt,
		CreatedAt:    eo.CreatedAt,
		SentAt:       eo.SentAt,
	}

	return res
}

func (eo InboxEvent) ToEntity() contracts.InboxEvent {
	res := contracts.InboxEvent{
		ID:           eo.ID,
		Topic:        eo.Topic,
		EventType:    eo.EventType,
		EventVersion: uint(eo.EventVersion),
		Key:          eo.Key,
		Payload:      eo.Payload,
		Status:       eo.Status,
		Attempts:     uint(eo.Attempts),
		NextRetryAt:  eo.NextRetryAt,
		CreatedAt:    eo.CreatedAt,
		ProcessedAt:  eo.ProcessedAt,
	}

	return res
}
