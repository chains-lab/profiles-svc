package handlers

import (
	"context"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
	"github.com/google/uuid"
)

func (s Service) UpdateOwnBirthday(ctx context.Context, req *svc.UpdateOwnBirthdayRequest) (*svc.Biography, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		Log(ctx, requestID).WithError(err).Errorf("invalid date format for birthday: %s", req.Birthday)

		return nil, responses.BadRequestError(ctx, requestID, responses.Violation{
			Field:       "birthday",
			Description: "invalid date format, expected YYYY-MM-DD",
		})
	}

	res, err := s.app.UpdateBirthday(ctx, meta.InitiatorID, birthday)
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to update user birthday")

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.Biography(res), nil
}
