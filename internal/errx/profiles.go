package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorProfileForUserDoesNotExist = ape.DeclareError("PROFILE_FOR_USER_DOES_NOT_EXIST")

var ErrorProfileForUserAlreadyExists = ape.DeclareError("PROFILE_FOR_USER_ALREADY_EXISTS")

var ErrorUsernameAlreadyTaken = ape.DeclareError("USERNAME_ALREADY_TAKEN")

var ErrorUsernameIsNotValid = ape.DeclareError("USERNAME_IS_NOT_VALID")

var ErrorSexIsNotValid = ape.DeclareError("SEX_IS_NOT_VALID")

var ErrorBirthdateIsNotValid = ape.DeclareError("BIRTHDATE_IS_NOT_VALID")

var ErrorUserTooYoung = ape.DeclareError("USER_TOO_YOUNG")
