package service

import (
	"context"
	"fmt"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
)

func (s Service) CreateOwnCabinet(ctx context.Context, req *svc.CreateCabinetRequest) (*svc.Cabinet, error) {
	meta := Meta(ctx)

	if meta.Role != roles.User {
		Log(ctx, meta.RequestID).Warnf(fmt.Sprintf(
			"User %s with role %s tried to create a cabinet, but only users can create cabinets",
			meta.InitiatorID, meta.Role),
		)

		return nil, responses.AppError(ctx, meta.RequestID, ape.ErrorOnlyUserCanHaveCabinet(
			fmt.Errorf("user %s with role %s tried to create a cabinet", meta.InitiatorID, meta.Role)))
	}

	cabinet, err := s.app.CreateCabinet(ctx, meta.InitiatorID, app.CreateCabinetInput{
		Username:    req.Username,
		Pseudonym:   req.Pseudonym,
		Description: req.Description,
		Avatar:      req.Avatar,
	})
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to create cabinet")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Cabinet(cabinet), nil
}
