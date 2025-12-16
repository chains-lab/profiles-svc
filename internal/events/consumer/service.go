package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/events/consumer/subscriber"
	"github.com/chains-lab/profiles-svc/internal/events/contracts"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
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

// U need to remake callbacks, they must don't use domain, only save to inbox (they can use domain but not in our case)
// and make 2 separate function to run - inbox reading and processing and kafka subscribing for topics

type Callbacks interface {
	UpdateEmployee(ctx context.Context, event kafka.Message) error
	UpdateCityAdmin(ctx context.Context, event kafka.Message) error
}

func Run(ctx context.Context, log logium.Logger, addr string, cb Callbacks) {
	var wg sync.WaitGroup

	accountSub := subscriber.New(addr, contracts.AccountsTopicV1, contracts.GroupProfilesSvc)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := accountSub.Subscribe(ctx, "employee.update", cb.UpdateEmployee); err != nil {
			log.Printf("employee listener stopped: %v", err)
		}
	}()

	<-ctx.Done()
	wg.Wait()
}
