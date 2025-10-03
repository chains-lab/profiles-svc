package enum

import "fmt"

const (
	SexMale   = "male"
	SexFemale = "female"
	SexOther  = "other"
)

var sexes = []string{
	SexMale,
	SexFemale,
	SexOther,
}

var ErrorSexNotSupported = fmt.Errorf("sex is not supported, must be one of: %v", GetAllSexes())

func CheckUserSexes(sex string) error {
	for _, s := range sexes {
		if s == sex {
			return nil
		}
	}

	return fmt.Errorf("%s: %w", sex, ErrorSexNotSupported)
}

func GetAllSexes() []string {
	return sexes
}
