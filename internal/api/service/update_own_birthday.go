package service

import (
	"context"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) UpdateOwnBirthday(ctx context.Context, req *svc.UpdateOwnBirthdayRequest) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Errorf("invalid date format for birthday: %s", req.Birthday)

		return nil, responses.BadRequestError(ctx, meta.RequestID, ape.Violation{
			Field:       "birthday",
			Description: "invalid date format, expected YYYY-MM-DD",
		})
	}

	err = s.app.UpdateBirthday(ctx, meta.InitiatorID, birthday)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update user birthday")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Infof("birthday for user %s has been updated to %s", meta.InitiatorID, req.Birthday)

	return &emptypb.Empty{}, nil
}
