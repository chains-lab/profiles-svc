package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) UpdateOwnSex(ctx context.Context, req *svc.UpdateOwnSexRequest) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	err := s.app.UpdateSex(ctx, meta.InitiatorID, req.Sex)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update user sex")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Infof("update sex for user %s to %s", meta.InitiatorID, req.Sex)

	return &emptypb.Empty{}, nil
}
