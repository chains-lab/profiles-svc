package ape

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ServiceName = "elector-cab-svc"

var (
	ErrorInternal = &Error{reason: ReasonInternal}

	ErrorOnlyUserCanHaveProfile      = &Error{reason: ReasonOnlyUserCanHaveProfile}
	ErrorProfileForUserDoesNotExist  = &Error{reason: ReasonProfileForUserDoesNotExist}
	ErrorProfileForUserAlreadyExists = &Error{reason: ReasonProfileForUserAlreadyExists}

	ErrorUsernameAlreadyTaken   = &Error{reason: ReasonUsernameAlreadyTaken}
	ErrorUsernameIsNotValid     = &Error{reason: ReasonUsernameIsNotValid}
	ErrorUsernameUpdateCooldown = &Error{reason: ReasonUsernameUpdateCooldown}

	ErrorBirthdayIsNotValid      = &Error{reason: ReasonBirthdayIsNotValid}
	ErrorBirthdayIsAlreadySet    = &Error{reason: ReasonBirthdayIsAlreadySet}
	ErrorSexIsNotValid           = &Error{reason: ReasonSexIsNotValid}
	ErrorSexUpdateCooldown       = &Error{reason: ReasonSexUpdateCooldown}
	ErrorResidenceIsNotValid     = &Error{reason: ReasonResidenceIsNotValid}
	ErrorResidenceUpdateCooldown = &Error{reason: ReasonResidenceUpdateCooldown}
)

func RaiseInternal(cause error) error {
	return &Error{
		code:    ErrorInternal.code,
		reason:  ErrorInternal.reason,
		message: "unexpected internal error occurred",
		cause:   cause,
	}
}

func RaiseProfileForUserDoesNotExist(cause error, user string) error {
	return &Error{
		code:    codes.NotFound,
		reason:  ErrorProfileForUserDoesNotExist.reason,
		message: fmt.Sprintf("profile for user %s does not exist", user),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ErrorProfileForUserDoesNotExist.reason,
				Domain: ServiceName,
			},
		},
	}
}

func RaiseProfileForUserAlreadyExists(cause error, user string) error {
	return &Error{
		code:    codes.AlreadyExists,
		reason:  ErrorProfileForUserAlreadyExists.reason,
		message: fmt.Sprintf("cabinet for user %s already exists", user),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ErrorProfileForUserAlreadyExists.reason,
				Domain: ServiceName,
			},
		},
	}
}

func RaiseOnlyUserCanHaveCabinetAndProfile(cause error) error {
	return &Error{
		code:    codes.PermissionDenied,
		reason:  ErrorOnlyUserCanHaveProfile.reason,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ErrorOnlyUserCanHaveProfile.reason,
				Domain: ServiceName,
			},
		},
	}
}

func RaiseUsernameAlreadyTaken(cause error, username string) error {
	return &Error{
		code:    codes.FailedPrecondition,
		reason:  ErrorUsernameAlreadyTaken.reason,
		message: fmt.Sprintf("username %s is already taken", username),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ErrorUsernameAlreadyTaken.reason,
				Domain: ServiceName,
			},
			&errdetails.ResourceInfo{
				ResourceType: "username",
				ResourceName: username,
				Description:  "This username is already in use by another account",
			},
		},
	}
}

func RaiseUsernameIsNotValid(cause error) error {
	return &Error{
		code:    codes.InvalidArgument,
		reason:  ErrorUsernameIsNotValid.reason,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ErrorUsernameIsNotValid.reason,
				Domain: ServiceName,
			},
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       "username",
				Description: "username is not valid, it must be 3-32 characters long, allowed characters are: a-z, A-Z, 0-9, _ (underscore), - (dash), . (dot)",
			}}},
		},
	}
}

func RaiseUsernameUpdateCooldown(cause error) error {
	return &Error{
		code:    codes.FailedPrecondition,
		reason:  ErrorUsernameUpdateCooldown.reason,
		message: "username can be updated only once per 14 days",
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ErrorUsernameUpdateCooldown.reason,
				Domain: ServiceName,
			},
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        "username_update_cooldown",
				Subject:     "username",
				Description: "username can be updated only once per 14 days",
			}}},
		},
	}
}

func RaiseBirthdayIsNotValid(cause error) error {
	return &Error{
		code:    codes.InvalidArgument,
		reason:  ReasonBirthdayIsNotValid,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ReasonBirthdayIsNotValid,
				Domain: ServiceName,
			},
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       "birthday",
				Description: "birthday is not valid, it must be in the past, but not more than 1900-01-01",
			}}},
		},
	}
}

func RaiseBirthdayIsAlreadySet(cause error) error {
	return &Error{
		code:    codes.FailedPrecondition,
		reason:  ReasonBirthdayIsAlreadySet,
		message: "birthday is already set",
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ReasonBirthdayIsAlreadySet,
				Domain: ServiceName,
			},
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        "birthday_already_set",
				Subject:     "birthday",
				Description: "birthday is already set and cannot be changed",
			}}},
		},
	}
}

func RaiseSexIsNotValid(cause error) error {
	return &Error{
		code:    codes.InvalidArgument,
		reason:  ReasonSexIsNotValid,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ReasonSexIsNotValid,
				Domain: ServiceName,
			},
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       "sex",
				Description: "sex is not valid", //TODO: add more details about valid values
			}}},
		},
	}
}

func RaiseSexUpdateCooldown(cause error) error {
	return &Error{
		code:    codes.FailedPrecondition,
		reason:  ReasonSexUpdateCooldown,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ReasonSexUpdateCooldown,
				Domain: ServiceName,
			},
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        "sex_update_cooldown",
				Subject:     "sex",
				Description: "sex can be updated only once per year",
			}}},
		},
	}
}

func RaiseResidenceIsNotValid(cause error) error {
	return &Error{
		code:    codes.InvalidArgument,
		reason:  ReasonResidenceIsNotValid,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ReasonResidenceIsNotValid,
				Domain: ServiceName,
			},
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       "residence",
				Description: "residence is not valid, it must be a valid country name",
			}}},
		},
	}
}

func RaiseResidenceUpdateCooldown(cause error) error {
	return &Error{
		code:    codes.FailedPrecondition,
		reason:  ReasonResidenceUpdateCooldown,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: ReasonResidenceUpdateCooldown,
				Domain: ServiceName,
			},
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        "residence_update_cooldown",
				Subject:     "residence",
				Description: "residence can be updated only once per 100 days",
			}}},
		},
	}
}
