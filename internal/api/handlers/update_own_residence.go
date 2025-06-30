package handlers

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
	"github.com/google/uuid"
)

func (s Service) UpdateOwnResidence(ctx context.Context, req *svc.UpdateOwnResidenceRequest) (*svc.Biography, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	bio, err := s.app.UpdateResidence(ctx, meta.InitiatorID, req.Country, req.City)
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to update user residence")

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.Biography(bio), nil
}
