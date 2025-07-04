package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) UpdateOwnResidence(ctx context.Context, req *svc.UpdateOwnResidenceRequest) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	err := s.app.UpdateResidence(ctx, meta.InitiatorID, app.UpdateResidenceInput{
		Country: req.Country,
		Region:  req.Region,
		City:    req.City,
	})
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update user residence")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Infof("residence for user %s has been updated to %s, %s, %s", meta.InitiatorID, req.Country, req.Region, req.City)

	return &emptypb.Empty{}, nil
}
