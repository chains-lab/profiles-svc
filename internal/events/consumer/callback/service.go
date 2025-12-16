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
		event contracts.InboxEvent,
	) error
}

type Service struct {
	domain Domain
	log    logium.Logger
}

func NewService(domain Domain) *Service {
	return &Service{
		domain: domain,
	}
}

//func decodeEnvelope[T any](b []byte) (events.Envelope[T], error) {
//	var env events.Envelope[T]
//	err := json.Unmarshal(b, &env)
//	return env, err
//}
