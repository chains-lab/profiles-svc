package callback

import (
	"context"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/events/contracts"
)

type Inbox interface {
	CreateInboxEvent(
		ctx context.Context,
		event contracts.InboxEvent,
	) error
}

type Service struct {
	inbox Inbox
	log   logium.Logger
}

func NewService(log logium.Logger, inbox Inbox) *Service {
	return &Service{
		inbox: inbox,
		log:   log,
	}
}
