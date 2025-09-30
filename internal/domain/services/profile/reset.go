package profile

import (
	"context"
	"crypto/rand"
	"fmt"

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

func (s Service) ResetUsername(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	usrnm, err := generateUsername()
	if err != nil {
		return models.Profile{}, err
	}

	u, err := s.UpdateUsername(ctx, userID, usrnm)
	if err != nil {
		return models.Profile{}, err
	}

	return u, nil
}

func (s Service) ResetUserProfile(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	empty := ""
	dmInput := Update{}
	dmInput.Pseudonym = &empty
	dmInput.Description = &empty
	dmInput.Avatar = &empty

	res, err := s.Update(ctx, userID, dmInput)
	if err != nil {
		return models.Profile{}, err
	}

	return res, nil
}
