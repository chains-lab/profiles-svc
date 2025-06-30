package ape

import "fmt"

var (
	ErrInternal                   = &BusinessError{reason: ReasonInternal}
	ErrPropertyUpdateNotAllowed   = &BusinessError{reason: ReasonPropertyUpdateNotAllowed}
	ErrPropertyIsNotValid         = &BusinessError{reason: ReasonPropertyIsNotValid}
	ErrUsernameAlreadyTaken       = &BusinessError{reason: ReasonUsernameAlreadyTaken}
	ErrCabinetForUserDoesNotExist = &BusinessError{reason: ReasonCabinetForUserDoesNotExist}
)

func ErrorInternal(cause error) error {
	return &BusinessError{
		reason:  ErrInternal.reason,
		message: "unexpected internal error occurred",
		cause:   cause,
	}
}

func ErrorUsernameAlreadyTaken(cause error, username string) error {
	return &BusinessError{
		reason:  ErrUsernameAlreadyTaken.reason,
		message: fmt.Sprintf("username %s is already taken", username),
		cause:   cause,
	}
}

func ErrorCabinetForUserDoesNotExist(cause error, user string) error {
	return &BusinessError{
		reason:  ErrCabinetForUserDoesNotExist.reason,
		message: fmt.Sprintf("cabinet for user %s does not exist", user),
		cause:   cause,
	}
}

func ErrorCabinetForUserAlreadyExists(cause error, user string) error {
	return &BusinessError{
		reason:  ErrCabinetForUserDoesNotExist.reason,
		message: fmt.Sprintf("cabinet for user %s already exists", user),
		cause:   cause,
	}
}

func ErrorPropertyUpdateNotAllowed(cause error) error {
	return &BusinessError{
		reason:  ErrPropertyUpdateNotAllowed.reason,
		message: cause.Error(),
		cause:   cause,
	}
}

func ErrorPropertyIsNotValid(cause error) error {
	return &BusinessError{
		reason:  ErrPropertyIsNotValid.reason,
		message: cause.Error(),
		cause:   cause,
	}
}
