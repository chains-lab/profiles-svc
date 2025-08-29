package constant

import "fmt"

const (
	UserSexMale   = "male"
	UserSexFemale = "female"
	UserSexOther  = "other"
)

var userSexes = []string{
	UserSexMale,
	UserSexFemale,
	UserSexOther,
}

var ErrorUserSexIsNotSupported = fmt.Errorf("user sex is not supported")

func ParseUserSex(sex string) error {
	for _, userSex := range userSexes {
		if userSex == sex {
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrorUserSexIsNotSupported, sex)
}

func GetAllUserSexes() []string {
	return userSexes
}
