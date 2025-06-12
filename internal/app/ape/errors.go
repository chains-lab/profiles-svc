package ape

type Error struct {
	Code    string
	Title   string
	Details string
	cause   error
}

func (e *Error) Unwrap() error {
	return e.cause
}

func ErrorInternal(cause error) *Error {
	return &Error{
		Code:    CodeInternal,
		Title:   "Internal Server Error",
		Details: "Internal server error",
		cause:   cause,
	}
}
