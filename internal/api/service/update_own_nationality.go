package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) UpdateOwnNationality(ctx context.Context, req *svc.UpdateOwnNationalityRequest) (*svc.Biography, error) {
	meta := Meta(ctx)

	bio, err := s.app.UpdateNationality(ctx, meta.InitiatorID, req.Nationality)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update user nationality")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Biography(bio), nil
}
