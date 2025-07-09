package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/logger"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) UpdateOwnSex(ctx context.Context, req *svc.UpdateOwnSexRequest) (*svc.Biography, error) {
	meta := Meta(ctx)

	res, err := s.app.UpdateSex(ctx, meta.InitiatorID, req.Sex)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to update user sex")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("update sex for user %s to %s", meta.InitiatorID, req.Sex)

	return responses.Biography(res), nil
}
