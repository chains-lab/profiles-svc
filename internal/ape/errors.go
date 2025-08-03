package ape

import (
	"fmt"

	"github.com/chains-lab/apperr"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ServiceName = "profiles-svc"

var ErrorInternal = &apperr.ErrorObject{Reason: ReasonInternal, Code: codes.Internal}

func RaiseInternal(cause error) error {
	return &apperr.ErrorObject{
		Code:    ErrorInternal.Code,
		Reason:  ErrorInternal.Reason,
		Message: "unexpected internal error occurred",
		Cause:   cause,
	}
}

var ErrorProfileForUserDoesNotExist = &apperr.ErrorObject{Reason: ReasonProfileForUserNotFound, Code: codes.NotFound}

func RaiseProfileForUserNotFound(cause error, userID uuid.UUID) error {
	return &apperr.ErrorObject{
		Code:    ErrorProfileForUserDoesNotExist.Code,
		Reason:  ErrorProfileForUserDoesNotExist.Reason,
		Message: fmt.Sprintf("profile for user with user_id: %s not found", userID),
		Cause:   cause,
		Details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "profile",
				ResourceName: fmt.Sprintf("profile:user_id:%s", userID),
				Description:  fmt.Sprintf("profile for user with user_id: %s does not exist", userID),
			},
		},
	}
}

func RaiseProfileForUserNotFoundByUsername(cause error, username string) error {
	return &apperr.ErrorObject{
		Code:    ErrorProfileForUserDoesNotExist.Code,
		Reason:  ErrorProfileForUserDoesNotExist.Reason,
		Message: fmt.Sprintf("profile for user with username: %s not found", username),
		Cause:   cause,
		Details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "profile",
				ResourceName: fmt.Sprintf("profile:username:%s", username),
				Description:  fmt.Sprintf("profile for user with username %s does not exist", username),
			},
		},
	}
}

// TODO idk its error is useless, maybe I should remove it

var ErrorProfileForUserAlreadyExists = &apperr.ErrorObject{Reason: ReasonProfileForUserAlreadyExists, Code: codes.AlreadyExists}

func RaiseProfileForUserAlreadyExists(cause error, userID uuid.UUID) error {
	return &apperr.ErrorObject{
		Code:    ErrorProfileForUserAlreadyExists.Code,
		Reason:  ErrorProfileForUserAlreadyExists.Reason,
		Message: fmt.Sprintf("profile for user with user_id: %s already exists", userID),
		Cause:   cause,
		Details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "profile",
				ResourceName: fmt.Sprintf("profile:user_id:%s", userID),
				Owner:        fmt.Sprintf("user:id:%s", userID),
				Description:  fmt.Sprintf("profile for user with user_id: %s already exists", userID),
			},
		},
	}
}

var ErrorOnlyUserCanHaveProfile = &apperr.ErrorObject{Reason: ReasonOnlyUserCanHaveProfile, Code: codes.PermissionDenied}

func RaiseOnlyUserCanHaveProfile(cause error) error {
	return &apperr.ErrorObject{
		Code:    ErrorOnlyUserCanHaveProfile.Code,
		Reason:  ErrorOnlyUserCanHaveProfile.Reason,
		Message: "only users with role user can have a profile",
		Cause:   cause,
	}
}

var ErrorUsernameAlreadyTaken = &apperr.ErrorObject{Reason: ReasonUsernameAlreadyTaken, Code: codes.FailedPrecondition}

func RaiseUsernameAlreadyTaken(cause error, username string) error {
	return &apperr.ErrorObject{
		Code:    ErrorUsernameAlreadyTaken.Code,
		Reason:  ErrorUsernameAlreadyTaken.Reason,
		Message: fmt.Sprintf("username %s is already taken", username),
		Cause:   cause,
		Details: []protoadapt.MessageV1{
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        ErrorUsernameAlreadyTaken.Reason,
				Subject:     fmt.Sprintf("profile:username:%s/username", username),
				Description: fmt.Sprintf("username %s is already taken", username),
			}}},
		},
	}
}

var ErrorUsernameIsNotValid = &apperr.ErrorObject{Reason: ReasonUsernameIsNotValid, Code: codes.InvalidArgument}

func RaiseUsernameIsNotValid(cause error) error {
	return &apperr.ErrorObject{
		Code:    ErrorUsernameIsNotValid.Code,
		Reason:  ErrorUsernameIsNotValid.Reason,
		Message: cause.Error(),
		Cause:   cause,
		Details: []protoadapt.MessageV1{
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field: "username",
				Description: "username is not valid, it must be 3-32 characters long, allowed characters are:" +
					" a-z, A-Z, 0-9, _ (underscore), - (dash), . (dot)",
				Reason: ErrorUsernameIsNotValid.Reason,
			}}},
		},
	}
}

var ErrorUsernameUpdateCooldown = &apperr.ErrorObject{Reason: ReasonUsernameUpdateCooldown, Code: codes.FailedPrecondition}

func RaiseUsernameUpdateCooldown(cause error, userID uuid.UUID) error {
	return &apperr.ErrorObject{
		Code:    ErrorUsernameUpdateCooldown.Code,
		Reason:  ErrorUsernameUpdateCooldown.Reason,
		Message: "username can be updated only once per 14 days",
		Cause:   cause,
		Details: []protoadapt.MessageV1{
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        ErrorUsernameUpdateCooldown.Reason,
				Subject:     fmt.Sprintf("profile:user_id:%s/username", userID),
				Description: "username can be updated only once per 14 days",
			}}},
		},
	}
}
