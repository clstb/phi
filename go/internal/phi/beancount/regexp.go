package beancount

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	DateRE        = regexp.MustCompile(`\d\d\d\d-\d\d-\d\d`)
	TransactionRE = regexp.MustCompile(
		fmt.Sprintf(
			`^(%s) (txn|!|\*) "([^"]*)"(?:$| "([^"]*)")?$`,
			DateRE.String(),
		),
	)
	PostingRE = regexp.MustCompile(
		fmt.Sprintf(
			`^  (%s)(?:(?:$|\s*)%s)?(?:(?:$| ){%s})?(?:(?:$| )(@|@@) %s)?$`,
			AccountRE.String(),
			AmountRE.String(),
			AmountRE.String(),
			AmountRE.String(),
		),
	)
	AccountRE = regexp.MustCompile(
		fmt.Sprintf(`(?:%s):\S*`, strings.Join([]string{
			"Assets",
			"Equity",
			"Expenses",
			"Income",
			"Liabilities",
		}, "|")),
	)
	AmountRE   = regexp.MustCompile(`(-?\d+(?:\.\d+)?) ([A-Z0-9]+)`)
	MetadataRE = regexp.MustCompile(`^  ([a-z].*): (\S*)$`)
	OpenRE     = regexp.MustCompile(fmt.Sprintf(
		`^(%s)\s+open\s+(%s)$`,
		DateRE.String(),
		AccountRE.String(),
	))
	SpecialRunesRE = regexp.MustCompile(`[^(\w|\s)]`)
)
