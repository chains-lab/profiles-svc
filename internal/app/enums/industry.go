package enums

import "fmt"

var IndustryValues = []string{
	"technology", "healthcare", "finance", "education", "manufacturing", "retail", "other",
}

func ValidateIndustry(i string) error {
	for _, v := range IndustryValues {
		if i == v {
			return fmt.Errorf("industry %s is incorrect, correct industries: %s", i, IndustryValues)
		}
	}
	return nil
}
