package ape

import (
	"fmt"
)

const ReasonUsernameAlreadyTaken = "USERNAME_ALREADY_Taken"

func ErrorUsernameAlreadyTaken(cause error, username string) error {
	return &Error{
		Reason:  ReasonUsernameAlreadyTaken,
		Details: fmt.Sprintf("username %s already taken", username),
		cause:   cause,
	}
}

const ReasonCabinetForUserDoesNotExist = "CABINET_FOR_USER_DOES_NOT_EXIST"

func ErrorCabinetForUserDoesNotExist(cause error, user string) error {
	return &Error{
		Reason:  ReasonCabinetForUserDoesNotExist,
		Details: fmt.Sprintf("cabinet for user %s does not exist", user),
		cause:   cause,
	}
}

const ReasonCabinetForUserAlreadyExists = "CABINET_FOR_USER_ALREADY_EXISTS"

func ErrorCabinetForUserAlreadyExists(cause error, user string) error {
	return &Error{
		Reason:  ReasonCabinetForUserAlreadyExists,
		Details: fmt.Sprintf("cabinet for user %s already exists", user),
		cause:   cause,
	}
}
