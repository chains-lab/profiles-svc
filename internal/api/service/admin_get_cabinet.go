package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"github.com/google/uuid"
)

func (s Service) AdminGetCabinet(ctx context.Context, req *svc.AdminGetCabinetRequest) (*svc.Cabinet, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("invalid user ID format")

		return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	cabinet, err := s.app.GetCabinetByUserID(ctx, userID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to get cabinet")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Cabinet(cabinet), nil
}
