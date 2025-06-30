package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetOwnBiography(ctx context.Context, _ *emptypb.Empty) (*svc.Biography, error) {
	meta := Meta(ctx)

	bio, err := s.app.GetUserBiographyByUserID(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to get biography")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Biography(bio), nil
}
