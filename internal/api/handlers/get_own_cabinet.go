package handlers

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
	"github.com/google/uuid"
)

func (s Service) GetOwnCabinet(ctx context.Context, _ *svc.Empty) (*svc.Cabinet, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	cabinet, err := s.app.GetCabinetByUserID(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to get cabinet")

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.Cabinet(cabinet), nil
}
