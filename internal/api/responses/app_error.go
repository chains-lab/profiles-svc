package responses

import (
	"context"
	"errors"
	"time"

	"github.com/chains-lab/apperr"
	"github.com/chains-lab/profiles-svc/internal/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func AppError(ctx context.Context, requestID uuid.UUID, err error) error {
	var appErr *apperr.ErrorObject
	if errors.As(err, &appErr) {

		st := status.New(appErr.Code, appErr.Message)

		details := appErr.Details
		details = append(details, &errdetails.ErrorInfo{
			Reason: appErr.Reason,
			Domain: ape.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		})
		details = append(details, &errdetails.RequestInfo{
			RequestId: requestID.String(),
		})
		if err != nil {
			return st.Err()
		}
	}

	return InternalError(ctx, &requestID)
}
