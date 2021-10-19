package parser

import (
	"bufio"
	"fmt"
	"regexp"
)

var (
	OpenRE = regexp.MustCompile(fmt.Sprintf(
		`^(%s)\s+open\s+(%s)$`,
		DateRE.String(),
		AccountRE.String(),
	))
)

type Opens []Open

func (os Opens) ByTinkID() map[string]Open {
	m := make(map[string]Open)
	for _, o := range os {
		for _, f := range o.Metadata {
			if f.Key == "tink_id" {
				m[f.Value] = o
			}
		}
	}

	return m
}

type Open struct {
	Date     string
	Account  string
	Metadata []MetadataField
}

func parseOpen(
	s string,
	scanner *bufio.Scanner,
	ledger Ledger,
) Ledger {
	matches := OpenRE.FindStringSubmatch(s)

	open := Open{
		Date:    matches[1],
		Account: matches[2],
	}

	for scanner.Scan() {
		s := scanner.Text()

		switch {
		case MetadataFieldRE.MatchString(s):
			open.Metadata = append(
				open.Metadata,
				parseMetadata(s),
			)
		default:
			return append(
				append(ledger, open),
				parse(s, scanner, ledger)...,
			)
		}
	}

	return append(ledger, open)
}

func (o Open) String() string {
	s := fmt.Sprintf("%s open %s\n", o.Date, o.Account)

	for _, md := range o.Metadata {
		s += md.String()
	}
	return s
}
