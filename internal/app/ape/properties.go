package ape

const ReasonPropertyUpdateNotAllowed = "PROPERTY_UPDATE_NOT_ALLOWED"

func ErrorPropertyUpdateNotAllowed(cause error) error {
	return &Error{
		Reason:  ReasonPropertyUpdateNotAllowed,
		Details: cause.Error(),
		cause:   cause,
	}
}

const ReasonPropertyIsNotValid = "PROPERTY_IS_NOT_VALID"

func ErrorPropertyIsNotValid(cause error) error {
	return &Error{
		Reason:  ReasonPropertyIsNotValid,
		Details: cause.Error(),
		cause:   cause,
	}
}
