package handlers

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
	"github.com/google/uuid"
)

func (s Service) UpdateOwnNationality(ctx context.Context, req *svc.UpdateOwnNationalityRequest) (*svc.Biography, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	bio, err := s.app.UpdateNationality(ctx, meta.InitiatorID, req.Nationality)
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to update user nationality")

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.Biography(bio), nil
}
