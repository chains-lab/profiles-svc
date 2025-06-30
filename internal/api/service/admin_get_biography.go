package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) AdminGetBiography(ctx context.Context, req *svc.AdminGetBiographyRequest) (*svc.Biography, error) {
	meta := Meta(ctx)

	bio, err := s.app.GetUserBiographyByUserID(ctx, meta.RequestID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to get biography")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Biography(bio), nil
}
