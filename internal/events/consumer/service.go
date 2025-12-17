package consumer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/events/consumer/subscriber"
	"github.com/chains-lab/profiles-svc/internal/events/contracts"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	addr      string
	inbox     inbox
	callbacks Callbacks
	log       logium.Logger
}

type Callbacks interface {
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

func New(log logium.Logger, addr string, inbox inbox, callbacks Callbacks) *Service {
	return &Service{
		addr:      addr,
		inbox:     inbox,
		log:       log,
		callbacks: callbacks,
	}
}

const eventInboxRetryDelay = 1 * time.Minute

func (s Service) Run(ctx context.Context) {
	var wg sync.WaitGroup

	accountSub := subscriber.New(s.addr, contracts.AccountsTopicV1, contracts.GroupProfilesSvc)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := accountSub.Subscribe(ctx, "account.create", s.callbacks.CreateAccount); err != nil {
			log.Printf("create account listener stopped: %v", err)
		}
	}()

	usernameSub := subscriber.New(s.addr, contracts.AccountsTopicV1, contracts.GroupProfilesSvc)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := usernameSub.Subscribe(ctx, "account.username.change", s.callbacks.UpdateUsername); err != nil {
			log.Printf("update username listener stopped: %v", err)
		}
	}()

	s.log.Info("starting events consumer", "addr", s.addr)
}

func (s Service) InboxWorker(ctx context.Context) {
	var wg sync.WaitGroup

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

		<-ctx.Done()
		wg.Wait()
	}
}
