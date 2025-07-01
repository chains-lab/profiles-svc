package responses

import (
	"context"
	"errors"

	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AppError(ctx context.Context, requestID uuid.UUID, err error) error {
	var appErr *ape.BusinessError
	if errors.As(err, &appErr) {
		var code codes.Code
		switch appErr.Reason() {
		case ape.ReasonInternal:
			code = codes.Internal

		case ape.ReasonPropertyUpdateNotAllowed:
			// попытка изменить поле, которое нельзя менять
			// семантически — предусловие не выполнено
			code = codes.FailedPrecondition

		case ape.ReasonPropertyIsNotValid,
			ape.ReasonUsernameIsNotValid:
			code = codes.InvalidArgument

		case ape.ReasonUsernameAlreadyTaken:
			// ресурс с таким уникальным ключом уже существует
			code = codes.AlreadyExists

		case ape.ReasonCabinetForUserDoesNotExist:
			// не найден ресурс «кабинет» у данного пользователя
			code = codes.NotFound

		default:
			code = codes.Unknown
		}

		st := status.New(code, appErr.Error())

		info := &errdetails.ErrorInfo{
			Reason: appErr.Reason(),
			Domain: "elector-cab-svc",
			Metadata: map[string]string{
				"request_id": requestID.String(),
			},
		}

		if code == codes.InvalidArgument {
			var fb []*errdetails.BadRequest_FieldViolation

			for _, v := range appErr.Violations() {
				fb = append(fb, &errdetails.BadRequest_FieldViolation{
					Field:       v.Field,
					Description: v.Description,
				})
			}
			br := &errdetails.BadRequest{FieldViolations: fb}

			st, err := st.WithDetails(info, br)
			if err != nil {
				return st.Err()
			}
		}

		return st.Err()
	}

	// всё, что «не BusinessError» — трактуем как Internal
	return status.Errorf(codes.Internal, "unexpected error")
}
