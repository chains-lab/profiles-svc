package service

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) GetCabinet(ctx context.Context, req *svc.GetCabinetRequest) (*svc.Cabinet, error) {
	meta := Meta(ctx)

	username := req.GetUsername()
	userID := req.GetUserId()

	var cabinet models.Cabinet
	var err error
	if username != "" {
		cabinet, err = s.app.GetCabinetByUserID(ctx, meta.InitiatorID)
		if err != nil {
			Log(ctx, meta.RequestID).WithError(err).Error("failed to get cabinet")

			return nil, responses.AppError(ctx, meta.RequestID, err)
		}
	} else if userID != "" {
		cabinet, err = s.app.GetCabinetByUsername(ctx, userID)
		if err != nil {
			Log(ctx, meta.RequestID).WithError(err).Error("failed to get cabinet by username")

			return nil, responses.AppError(ctx, meta.RequestID, err)
		}
	} else {
		return nil, responses.BadRequestError(ctx, meta.RequestID, ape.Violation{
			Field:       "username",
			Description: "either username is required or user_id",
		}, ape.Violation{
			Field:       "user_id",
			Description: "either user_id is required or username",
		})
	}

	return responses.Cabinet(cabinet), nil
}
