package beanacount

import (
	"fmt"
	"strings"
)

func (a *AccountType) String() string {
	name := toCamelCase(a.Name)
	if a.ID != "" {
		return fmt.Sprintf("%s:%s:%s", a.FinancialInstitutionId, a.ID, name)
	} else {
		return fmt.Sprintf("%s:%s", a.FinancialInstitutionId, name)
	}
}

func toCamelCase(input string) (camelCase string) {

	isToUpper := false

	for k, v := range input {
		if k == 0 {
			camelCase = strings.ToUpper(string(input[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == ' ' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return

}
