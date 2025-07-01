package service

import (
	"context"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	"github.com/chains-lab/elector-cab-svc/internal/app"
	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/electorcab"
	"github.com/google/uuid"
)

func (s Service) AdminUpdateBiography(ctx context.Context, req *svc.UpdateBiographyByAdminRequest) (*svc.Biography, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("invalid user ID format")

		return nil, responses.BadRequestError(ctx, meta.RequestID, ape.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	var birthday time.Time
	if req.Birthday != nil {
		birthday, err = time.Parse("02-01-2006", *req.Birthday)

		if err != nil {
			Log(ctx, meta.RequestID).WithError(err).Error("invalid birthday format")

			return nil, responses.BadRequestError(ctx, meta.RequestID, ape.Violation{
				Field:       "birthday",
				Description: "invalid date format, expected DD-MM-YYYY",
			})
		}
	}

	bio, err := s.app.AdminUpdateBiography(ctx, userID, app.UpdateBiographyInput{
		Birthday:        &birthday,
		Sex:             req.Sex,
		City:            req.City,
		Region:          req.Region,
		Country:         req.Country,
		Nationality:     req.Nationality,
		PrimaryLanguage: req.PrimaryLanguage,
	})
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to update biography")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Biography(bio), nil
}
