package callback

import (
	"context"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/domain/entity"
	"github.com/chains-lab/profiles-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type Domain interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, username string) (entity.Profile, error)
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (entity.Profile, error)
}

type Inbox interface {
	CreateInboxEvent(
		ctx context.Context,
		event contracts.Message,
	) error
}

type Service struct {
	domain Domain
	inbox  Inbox
	log    logium.Logger
}

func NewService(log logium.Logger, inbox Inbox, domain Domain) *Service {
	return &Service{
		domain: domain,
		inbox:  inbox,
		log:    log,
	}
}
