package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetOwnCabinet(ctx context.Context, _ *emptypb.Empty) (*svc.Cabinet, error) {
	meta := Meta(ctx)

	cabinet, err := s.app.GetCabinetByUserID(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to get cabinet")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Cabinet(cabinet), nil
}
