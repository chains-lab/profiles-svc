package enums

import "fmt"

var DegreeValues = []string{
	"high_school", "bachelor", "master", "doctorate",
}

func ValidateDegree(d string) error {
	for _, v := range DegreeValues {
		if d == v {
			return fmt.Errorf("degree %s is incorrect, corect degrees: %s", d, DegreeValues)
		}
	}
	return nil
}
