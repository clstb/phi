package parser

import (
	"fmt"
	"regexp"

	"github.com/shopspring/decimal"
)

var (
	PostingRE = regexp.MustCompile(
		fmt.Sprintf(
			`^  (%s)(?:(?:$|\s*)%s)?(?:(?:$| ){%s})?(?:(?:$| )(?:@|@@) %s)?$`,
			AccountRE.String(),
			AmountRE.String(),
			AmountRE.String(),
			AmountRE.String(),
		),
	)
)

type Posting struct {
	Account string
	Units   Amount
	Cost    Amount
	Price   Amount
}

func parsePosting(
	s string,
) Posting {
	matches := PostingRE.FindStringSubmatch(s)

	posting := Posting{
		Account: matches[1],
	}

	if len(matches[2]) != 0 {
		units, _ := decimal.NewFromString(matches[2])
		posting.Units = Amount{
			Decimal:  units,
			Currency: matches[3],
		}
	}

	if len(matches[4]) != 0 {
		cost, _ := decimal.NewFromString(matches[4])
		posting.Cost = Amount{
			Decimal:  cost,
			Currency: matches[5],
		}

	}

	if len(matches[6]) != 0 {
		price, _ := decimal.NewFromString(matches[6])
		posting.Price = Amount{
			Decimal:  price,
			Currency: matches[6],
		}
	}

	return posting
}

func (p Posting) String() string {
	if p.Units.IsZero() {
		return fmt.Sprintf("  %s\n", p.Account)
	}

	s := fmt.Sprintf(
		"  %s %s",
		p.Account,
		p.Units,
	)
	if !p.Cost.IsZero() {
		s += fmt.Sprintf(" {%s}", p.Cost)
	}
	if !p.Price.IsZero() {
		s += fmt.Sprintf(" @ %s", p.Price)
	}
	s += "\n"

	return s
}

func (p Posting) Weight() Amount {
	if !p.Cost.IsZero() {
		return p.Units.Mul(p.Cost)
	}
	if !p.Price.IsZero() {
		return p.Units.Mul(p.Price)
	}

	return p.Units
}
