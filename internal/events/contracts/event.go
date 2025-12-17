package contracts

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Topic        string
	EventType    string
	EventVersion uint
	Key          string
	Payload      json.RawMessage
}

type OutboxEvent struct {
	ID           uuid.UUID
	Topic        string
	EventType    string
	EventVersion uint
	Key          string
	Payload      json.RawMessage
	Status       string
	Attempts     uint
	NextRetryAt  time.Time
	CreatedAt    time.Time
	SentAt       *time.Time
}

func (e OutboxEvent) ToMessage() Message {
	return Message{
		Topic:        e.Topic,
		EventType:    e.EventType,
		EventVersion: e.EventVersion,
		Key:          e.Key,
		Payload:      e.Payload,
	}
}

type InboxEvent struct {
	ID           uuid.UUID
	Topic        string
	EventType    string
	EventVersion uint
	Key          string
	Payload      json.RawMessage
	Status       string
	Attempts     uint
	NextRetryAt  time.Time
	CreatedAt    time.Time
	ProcessedAt  *time.Time
}

func (e InboxEvent) ToMessage() Message {
	return Message{
		Topic:        e.Topic,
		EventType:    e.EventType,
		EventVersion: e.EventVersion,
		Key:          e.Key,
		Payload:      e.Payload,
	}
}

func (e InboxEvent) IsNil() bool {
	return e.ID == uuid.Nil
}
