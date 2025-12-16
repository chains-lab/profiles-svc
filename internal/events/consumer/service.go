package consumer

import (
	"context"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type Service struct {
	addr  string
	inbox inbox
	log   logium.Logger
}

type inbox interface {
	CreateInboxEvent(
		ctx context.Context,
		event contracts.InboxEvent,
	) error

	GetPendingInboxEvents(
		ctx context.Context,
		limit int32,
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

func New(log logium.Logger, addr string, inbox inbox) *Service {
	return &Service{
		addr:  addr,
		inbox: inbox,
		log:   log,
	}
}

const eventInboxRetryDelay = 1 * time.Minute

func (s *Service) Run(ctx context.Context) error {
	s.log.Info("starting events consumer", "addr", s.addr)

	for {
		events, err := s.inbox.GetPendingInboxEvents(ctx, 10)
		if err != nil {
			s.log.Error("failed to get pending inbox events", "error", err)
			time.Sleep(eventInboxRetryDelay)
			continue
		}

		if len(events) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		var processedIDs []uuid.UUID
		for _, event := range events {
			s.log.Info("processing inbox event", "id", event.ID, "type", event.EventType)
		}
		if len(processedIDs) > 0 {
			err = s.inbox.MarkInboxEventsAsProcessed(ctx, processedIDs)
			if err != nil {
				s.log.Error("failed to mark inbox events as processed", "error", err)
				time.Sleep(eventInboxRetryDelay)
				continue
			}
		}
	}
}
