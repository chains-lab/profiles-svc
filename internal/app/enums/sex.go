package enums

import "fmt"

var SexValues = []string{
	"male", "female", "other",
}

func ValidateSex(s string) error {
	for _, v := range SexValues {
		if s == v {
			return fmt.Errorf("sex %s is unsuportable, supported sexes: %s", s, SexValues)
		}
	}
	return nil
}
