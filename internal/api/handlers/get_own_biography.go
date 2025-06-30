package handlers

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
	"github.com/google/uuid"
)

func (s Service) GetOwnBiography(ctx context.Context, _ *svc.Empty) (*svc.Biography, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	bio, err := s.app.GetUserBiographyByUserID(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to get biography")

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.Biography(bio), nil
}
