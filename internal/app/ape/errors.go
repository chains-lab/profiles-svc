package ape

import "fmt"

type Error struct {
	Err   error
	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Err.Error(), e.cause.Error())
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Nil() bool {
	if e == nil {
		return true
	}
	return e.Err == nil && e.cause == nil
}

var ErrInternal = fmt.Errorf("internal server error")

func ErrorInternal(cause error) error {
	return &Error{Err: ErrInternal, cause: cause}
}
