package ape

const (
	ReasonBadRequest   = "BAD_REQUEST"
	ReasonInternal     = "INTERNAL_ERROR"
	ReasonUnauthorized = "UNAUTHORIZED"

	ReasonOnlyUserCanHaveProfile      = "ONLY_USER_CAN_HAVE_PROFILE"
	ReasonProfileForUserNotFound      = "PROFILE_FOR_USER_NOT_FOUND"
	ReasonProfileForUserAlreadyExists = "PROFILE_FOR_USER_ALREADY_EXISTS"

	ReasonUsernameAlreadyTaken   = "USERNAME_ALREADY_TAKEN"
	ReasonUsernameIsNotValid     = "USERNAME_IS_NOT_VALID"
	ReasonUsernameUpdateCooldown = "USERNAME_UPDATE_COOLDOWN"
)
