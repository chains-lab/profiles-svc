package enums

import "fmt"

var IncomeValues = []string{
	"below_30k", "30k_to_50k", "50k_to_70k", "70k_to_100k", "above_100k",
}

func ValidateIncome(i string) error {
	for _, v := range IncomeValues {
		if i == v {
			return fmt.Errorf("income %s is incorrect, correct incomes: %s", i, IncomeValues)
		}
	}
	return nil
}
