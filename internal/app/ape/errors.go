package ape

import "fmt"

type Error struct {
	Reason  string // similar to CODE in HTTP API errors
	Details error  // for internal use in application
	cause   error  // the original error that caused this error, if any
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Reason.Error(), e.cause.Error())
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Nil() bool {
	if e == nil {
		return true
	}
	return e.Details == nil && e.cause == nil
}

const ReasonUserDoesNotExist = "USER_DOES_NOT_EXIST"

var ErrInternal = fmt.Errorf("internal server error")

func ErrorInternal(cause error) error {
	return &Error{
		Reason:  ReasonUserDoesNotExist,
		Details: ErrInternal,
		cause:   cause,
	}
}
