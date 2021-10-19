package parser

import "regexp"

var (
	DateRE = regexp.MustCompile(`\d\d\d\d-\d\d-\d\d`)
)
