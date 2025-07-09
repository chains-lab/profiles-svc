package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	"github.com/chains-lab/elector-cab-svc/internal/logger"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) UpdateOwnResidence(ctx context.Context, req *svc.UpdateOwnResidenceRequest) (*svc.Biography, error) {
	meta := Meta(ctx)

	res, err := s.app.UpdateResidence(ctx, meta.InitiatorID, app.UpdateResidenceInput{
		Country: req.Country,
		Region:  req.Region,
		City:    req.City,
	})
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to update user residence")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("residence for user %s has been updated to %s, %s, %s", meta.InitiatorID, req.Country, req.Region, req.City)

	return responses.Biography(res), nil
}
