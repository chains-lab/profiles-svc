package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) UpdateOwnSex(ctx context.Context, req *svc.UpdateOwnSexRequest) (*svc.Biography, error) {
	meta := Meta(ctx)

	bio, err := s.app.UpdateSex(ctx, meta.InitiatorID, req.Sex)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update user sex")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Biography(bio), nil
}
