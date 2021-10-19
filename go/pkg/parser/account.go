package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	AccountTypes = []string{
		"Assets",
		"Equity",
		"Expenses",
		"Income",
		"Liabilities",
	}
	AccountRE = regexp.MustCompile(
		fmt.Sprintf(`(?:%s):\S*`, strings.Join(AccountTypes, "|")),
	)
)
