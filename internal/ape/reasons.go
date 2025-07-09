package ape

const (
	ReasonBadRequest = "BAD_REQUEST"
	ReasonInternal   = "INTERNAL_ERROR"

	ReasonOnlyUserCanHaveProfile      = "ONLY_USER_CAN_HAVE_PROFILE"
	ReasonProfileForUserDoesNotExist  = "PROFILE_FOR_USER_DOES_NOT_EXIST"
	ReasonProfileForUserAlreadyExists = "PROFILE_FOR_USER_ALREADY_EXISTS"

	ReasonUsernameAlreadyTaken    = "USERNAME_ALREADY_TAKEN"
	ReasonUsernameIsNotValid      = "USERNAME_IS_NOT_VALID"
	ReasonUsernameUpdateCooldown  = "USERNAME_UPDATE_COOLDOWN"
	ReasonBirthdayIsNotValid      = "BIRTHDAY_IS_NOT_VALID"
	ReasonBirthdayIsAlreadySet    = "BIRTHDAY_IS_ALREADY_SET"
	ReasonSexIsNotValid           = "SEX_IS_NOT_VALID"
	ReasonSexUpdateCooldown       = "SEX_UPDATE_COOLDOWN"
	ReasonResidenceIsNotValid     = "RESIDENCE_IS_NOT_VALID"
	ReasonResidenceUpdateCooldown = "RESIDENCE_UPDATE_COOLDOWN"
)
