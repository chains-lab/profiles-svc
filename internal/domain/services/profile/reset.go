package profile

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/chains-lab/profiles-svc/internal/domain/errx"
	"github.com/chains-lab/profiles-svc/internal/domain/models"
	"github.com/google/uuid"
)

func generateUsername() (string, error) {
	const (
		prefix = "citizen"
		digits = 8
	)
	buf := make([]byte, digits)
	if _, err := rand.Read(buf); err != nil {
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("cannot generate random digits: %w", err),
		)
	}
	for i := 0; i < digits; i++ {
		buf[i] = '0' + (buf[i] % 10)
	}
	return prefix + string(buf), nil
}

func (s Service) ResetProfile(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	p, err := s.GetByID(ctx, userID)
	if err != nil {
		return models.Profile{}, err
	}

	usrnm, err := generateUsername()
	if err != nil {
		return models.Profile{}, err
	}

	now := time.Now().UTC()

	err = s.db.ResetProfile(ctx, userID, usrnm, now)
	if err != nil {
		return models.Profile{}, err
	}

	p.Username = usrnm
	p.Avatar = nil
	p.Pseudonym = nil
	p.Description = nil
	p.UpdatedAt = now

	return p, nil
}
