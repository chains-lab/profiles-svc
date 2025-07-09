package references

import "fmt"

var sexEnums = []string{
	"man",
	"woman",
	"other",
}

func ValidateSex(s string) error {
	for _, v := range sexEnums {
		if s == v {
			return nil
		}
	}
	return fmt.Errorf("invalid sex value: %s, must be one of %v", s, sexEnums)
}
