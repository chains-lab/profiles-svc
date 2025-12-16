package pgdb

import (
	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/chains-lab/profiles-svc/internal/events/contracts"
)

func (p Profile) ToEntity() entity.Profile {
	profile := entity.Profile{
		AccountID: p.AccountID,
		Username:  p.Username,
		Official:  p.Official,

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	if p.Pseudonym.Valid {
		profile.Pseudonym = &p.Pseudonym.String
	}
	if p.Description.Valid {
		profile.Description = &p.Description.String
	}
	if p.Avatar.Valid {
		profile.Avatar = &p.Avatar.String
	}

	return profile
}

func (eo OutboxEvent) ToEntity() contracts.OutboxEvent {
	res := contracts.OutboxEvent{
		ID:           eo.ID,
		Topic:        eo.Topic,
		EventType:    eo.EventType,
		EventVersion: eo.EventVersion,
		Key:          eo.Key,
		Payload:      eo.Payload,
		Status:       string(eo.Status),
		Attempts:     eo.Attempts,
		NextRetryAt:  eo.NextRetryAt,
		CreatedAt:    eo.CreatedAt,
	}
	if eo.SentAt.Valid {
		t := eo.SentAt.Time
		res.SentAt = &t
	}

	return res
}

func (eo InboxEvent) ToEntity() contracts.InboxEvent {
	res := contracts.InboxEvent{
		ID:           eo.ID,
		Topic:        eo.Topic,
		EventType:    eo.EventType,
		EventVersion: eo.EventVersion,
		Key:          eo.Key,
		Payload:      eo.Payload,
		Status:       string(eo.Status),
		Attempts:     eo.Attempts,
		NextRetryAt:  eo.NextRetryAt,
		CreatedAt:    eo.CreatedAt,
	}
	if eo.ProcessedAt.Valid {
		t := eo.ProcessedAt.Time
		res.ProcessedAt = &t
	}

	return res
}
