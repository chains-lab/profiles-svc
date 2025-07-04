package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) UpdateOwnNationality(ctx context.Context, req *svc.UpdateOwnNationalityRequest) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	err := s.app.UpdateNationality(ctx, meta.InitiatorID, req.Nationality)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update user nationality")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Infof("nationality for user %s has been updated to %s", meta.InitiatorID, req.Nationality)

	return &emptypb.Empty{}, nil
}
