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
			// внутренняя неожиданная ошибка
			code = codes.Internal

		case ape.ReasonPropertyUpdateNotAllowed:
			// попытка изменить поле, которое нельзя менять
			// семантически — предусловие не выполнено
			code = codes.FailedPrecondition

		case ape.ReasonPropertyIsNotValid:
			// в теле запроса пришло некорректное значение поля
			code = codes.InvalidArgument

		case ape.ReasonUsernameAlreadyTaken:
			// ресурс с таким уникальным ключом уже существует
			code = codes.AlreadyExists

		case ape.ReasonCabinetForUserDoesNotExist:
			// не найден ресурс «кабинет» у данного пользователя
			code = codes.NotFound

		default:
			// неожиданный бизнес-код
			code = codes.Unknown
		}

		st := status.New(code, appErr.Error())
		st, errWithDetails := st.WithDetails(
			&errdetails.ErrorInfo{
				Reason: appErr.Reason(),
				Domain: "elector-cab.yourdomain.com", // ваше API/сервис
				Metadata: map[string]string{
					"request_id": requestID.String(),
				},
			},
		)
		if errWithDetails != nil {
			return st.Err()
		}
		return st.Err()
	}

	// всё, что «не BusinessError» — трактуем как Internal
	return status.Errorf(codes.Internal, "unexpected error")
}
