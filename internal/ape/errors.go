package ape

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ServiceName = "elector-cab-svc"

var ErrorInternal = &Error{reason: ReasonInternal, code: codes.Internal}

func RaiseInternal(cause error) error {
	return &Error{
		code:    ErrorInternal.code,
		reason:  ErrorInternal.reason,
		message: "unexpected internal error occurred",
		cause:   cause,
	}
}

var ErrorProfileForUserDoesNotExist = &Error{reason: ReasonProfileForUserNotFound, code: codes.NotFound}

func RaiseProfileForUserNotFound(cause error, userID uuid.UUID) error {
	return &Error{
		code:    ErrorProfileForUserDoesNotExist.code,
		reason:  ErrorProfileForUserDoesNotExist.reason,
		message: fmt.Sprintf("profile for user with user_id: %s not found", userID),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "profile",
				ResourceName: fmt.Sprintf("profile:user_id:%s", userID),
				Description:  fmt.Sprintf("profile for user with user_id: %s does not exist", userID),
			},
		},
	}
}

func RaiseProfileForUserNotFoundByUsername(cause error, username string) error {
	return &Error{
		code:    ErrorProfileForUserDoesNotExist.code,
		reason:  ErrorProfileForUserDoesNotExist.reason,
		message: fmt.Sprintf("profile for user with username: %s not found", username),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "profile",
				ResourceName: fmt.Sprintf("profile:username:%s", username),
				Description:  fmt.Sprintf("profile for user with username %s does not exist", username),
			},
		},
	}
}

var ErrorProfileForUserAlreadyExists = &Error{reason: ReasonProfileForUserAlreadyExists, code: codes.AlreadyExists}

func RaiseProfileForUserAlreadyExists(cause error, userID uuid.UUID) error {
	return &Error{
		code:    ErrorProfileForUserAlreadyExists.code,
		reason:  ErrorProfileForUserAlreadyExists.reason,
		message: fmt.Sprintf("cabinet for user with user_id: %s already exists", userID),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "profile",
				ResourceName: fmt.Sprintf("profile:user_id:%s", userID),
				Owner:        fmt.Sprintf("user:id:%s", userID),
				Description:  fmt.Sprintf("cabinet for user with user_id: %s already exists", userID),
			},
		},
	}
}

var ErrorOnlyUserCanHaveProfile = &Error{reason: ReasonOnlyUserCanHaveProfile, code: codes.PermissionDenied}

func RaiseOnlyUserCanHaveCabinetAndProfile(cause error) error {
	return &Error{
		code:    ErrorOnlyUserCanHaveProfile.code,
		reason:  ErrorOnlyUserCanHaveProfile.reason,
		message: "only users with role user can have a profile and cabinet",
		cause:   cause,
	}
}

var ErrorUsernameAlreadyTaken = &Error{reason: ReasonUsernameAlreadyTaken, code: codes.FailedPrecondition}

func RaiseUsernameAlreadyTaken(cause error, username string) error {
	return &Error{
		code:    ErrorUsernameAlreadyTaken.code,
		reason:  ErrorUsernameAlreadyTaken.reason,
		message: fmt.Sprintf("username %s is already taken", username),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        ErrorUsernameAlreadyTaken.reason,
				Subject:     fmt.Sprintf("profile:username:%s/username", username),
				Description: fmt.Sprintf("username %s is already taken", username),
			}}},
		},
	}
}

var ErrorUsernameIsNotValid = &Error{reason: ReasonUsernameIsNotValid, code: codes.InvalidArgument}

func RaiseUsernameIsNotValid(cause error) error {
	return &Error{
		code:    ErrorUsernameIsNotValid.code,
		reason:  ErrorUsernameIsNotValid.reason,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field: "username",
				Description: "username is not valid, it must be 3-32 characters long, allowed characters are:" +
					" a-z, A-Z, 0-9, _ (underscore), - (dash), . (dot)",
				Reason: ErrorUsernameIsNotValid.reason,
			}}},
		},
	}
}

var ErrorUsernameUpdateCooldown = &Error{reason: ReasonUsernameUpdateCooldown, code: codes.FailedPrecondition}

func RaiseUsernameUpdateCooldown(cause error, userID uuid.UUID) error {
	return &Error{
		code:    ErrorUsernameUpdateCooldown.code,
		reason:  ErrorUsernameUpdateCooldown.reason,
		message: "username can be updated only once per 14 days",
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        ErrorUsernameUpdateCooldown.reason,
				Subject:     fmt.Sprintf("profile:user_id:%s/username", userID),
				Description: "username can be updated only once per 14 days",
			}}},
		},
	}
}

var ErrorBirthdayIsNotValid = &Error{reason: ReasonBirthdayIsNotValid, code: codes.InvalidArgument}

func RaiseBirthdayIsNotValid(cause error) error {
	return &Error{
		code:    ErrorBirthdayIsNotValid.code,
		reason:  ErrorBirthdayIsNotValid.reason,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       "birthday",
				Description: "birthday is not valid, it must be in the past, but not more than 1900-01-01",
				Reason:      ErrorBirthdayIsNotValid.reason,
			}}},
		},
	}
}

var ErrorBirthdayIsAlreadySet = &Error{reason: ReasonBirthdayIsAlreadySet, code: codes.FailedPrecondition}

func RaiseBirthdayIsAlreadySet(cause error, userID uuid.UUID) error {
	return &Error{
		code:    ErrorBirthdayIsAlreadySet.code,
		reason:  ErrorBirthdayIsAlreadySet.reason,
		message: "birthday is already set",
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        ErrorBirthdayIsAlreadySet.reason,
				Subject:     fmt.Sprintf("profile:user_id:%s/birthday", userID),
				Description: "birthday is already set and cannot be changed",
			}}},
		},
	}
}

var ErrorSexIsNotValid = &Error{reason: ReasonSexIsNotValid, code: codes.InvalidArgument}

func RaiseSexIsNotValid(cause error) error {
	return &Error{
		code:    ErrorSexIsNotValid.code,
		reason:  ErrorSexIsNotValid.reason,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       "sex",
				Reason:      ErrorSexIsNotValid.reason,
				Description: "sex is not valid",
			}}},
		},
	}
}

var ErrorSexUpdateCooldown = &Error{reason: ReasonSexUpdateCooldown, code: codes.FailedPrecondition}

func RaiseSexUpdateCooldown(cause error, userID uuid.UUID) error {
	return &Error{
		code:    ErrorSexUpdateCooldown.code,
		reason:  ErrorSexUpdateCooldown.reason,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        ErrorSexUpdateCooldown.reason,
				Subject:     fmt.Sprintf("profile:user_id:%s/sex", userID),
				Description: "sex can be updated only once per year",
			}}},
		},
	}
}

var ErrorResidenceIsNotValid = &Error{reason: ReasonResidenceIsNotValid, code: codes.InvalidArgument}

func RaiseResidenceIsNotValid(cause error) error {
	return &Error{
		code:    codes.InvalidArgument,
		reason:  ErrorResidenceIsNotValid.reason,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       "residence",
				Reason:      ErrorResidenceIsNotValid.reason,
				Description: "residence is not valid, it must be a valid country name",
			}}},
		},
	}
}

var ErrorResidenceUpdateCooldown = &Error{reason: ReasonResidenceUpdateCooldown, code: codes.FailedPrecondition}

func RaiseResidenceUpdateCooldown(cause error, userID uuid.UUID) error {
	return &Error{
		code:    ErrorResidenceUpdateCooldown.code,
		reason:  ErrorResidenceUpdateCooldown.reason,
		message: cause.Error(),
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        ErrorResidenceUpdateCooldown.reason,
				Subject:     fmt.Sprintf("profile:user_id:%s/residence", userID),
				Description: "residence can be updated only once per 100 days",
			}}},
		},
	}
}
