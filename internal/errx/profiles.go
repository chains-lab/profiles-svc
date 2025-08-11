package errx

import (
	"context"
	"fmt"

	"github.com/chains-lab/profiles-svc/internal/api/grpc/meta"
	"github.com/chains-lab/profiles-svc/internal/constant"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorProfileForUserDoesNotExist = ape.Declare("PROFILE_FOR_USER_DOES_NOT_EXIST")

func RaiseProfileForUserNotFound(ctx context.Context, cause error, userID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("profile for user %s not found", userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorProfileForUserDoesNotExist.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorProfileForUserDoesNotExist.Raise(cause, st)
}

func RaiseProfileForUserNotFoundByUsername(ctx context.Context, cause error, username string) error {
	st := status.New(codes.NotFound, fmt.Sprintf("profile for username %q not found", username))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorProfileForUserDoesNotExist.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorProfileForUserDoesNotExist.Raise(cause, st)
}

var ErrorProfileForUserAlreadyExists = ape.Declare("PROFILE_FOR_USER_ALREADY_EXISTS")

func RaiseProfileForUserAlreadyExists(ctx context.Context, cause error, userID uuid.UUID) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("profile for user %s already exists", userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorProfileForUserAlreadyExists.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorProfileForUserAlreadyExists.Raise(cause, st)
}

var ErrorOnlyUserCanHaveProfile = ape.Declare("ONLY_USER_CAN_HAVE_PROFILE")

func RaiseOnlyUserCanHaveProfile(ctx context.Context, cause error) error {
	st := status.New(codes.PermissionDenied, "only users with role 'user' can have a profile")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorOnlyUserCanHaveProfile.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorOnlyUserCanHaveProfile.Raise(cause, st)
}

var ErrorUsernameAlreadyTaken = ape.Declare("USERNAME_ALREADY_TAKEN")

func RaiseUsernameAlreadyTaken(ctx context.Context, cause error, username string) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("username %q is already taken", username))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUsernameAlreadyTaken.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorUsernameAlreadyTaken.Raise(cause, st)
}

var ErrorUsernameIsNotValid = ape.Declare("USERNAME_IS_NOT_VALID")

func RaiseUsernameIsNotValid(ctx context.Context, cause error) error {
	// message из cause — ок, но без раскрытия лишних деталей
	st := status.New(codes.InvalidArgument, cause.Error())
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUsernameIsNotValid.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorUsernameIsNotValid.Raise(cause, st)
}

var ErrorUsernameUpdateCooldown = ape.Declare("USERNAME_UPDATE_COOLDOWN")

func RaiseUsernameUpdateCooldown(ctx context.Context, cause error, userID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, "username can be updated only once per 14 days")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUsernameUpdateCooldown.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorUsernameUpdateCooldown.Raise(cause, st)
}
