package ape

import (
	"fmt"
)

type Error struct {
	Reason  string // similar to CODE in HTTP API errors
	Details string // additional details about the error
	cause   error  // the original error that caused this error, if any
}

func (e *Error) Error() string {
	return fmt.Sprintf("reason %s: | details: %s | cause: %v", e.Reason, e.Details, e.cause.Error())
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Nil() bool {
	if e == nil {
		return true
	}
	return e.cause == nil
}

const ReasonErrorInternal = "INTERNAL_ERROR"

func ErrorInternal(cause error) error {
	return &Error{
		Reason:  ReasonErrorInternal,
		Details: "internal server error",
		cause:   cause,
	}
}
