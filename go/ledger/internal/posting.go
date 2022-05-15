package internal

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func parsePosting(
	s string,
) Posting {
	matches := PostingRE.FindStringSubmatch(s)

	posting := Posting{
		Account: matches[1],
	}

	if len(matches[2]) != 0 {
		units, _ := decimal.NewFromString(matches[2])
		posting.Units = AmountType{
			Decimal:  units,
			Currency: matches[3],
		}
	}

	if len(matches[4]) != 0 {
		cost, _ := decimal.NewFromString(matches[4])
		posting.Cost = AmountType{
			Decimal:  cost,
			Currency: matches[5],
		}
	}

	if len(matches[6]) != 0 {
		price, _ := decimal.NewFromString(matches[7])
		posting.Price = AmountType{
			Decimal:  price,
			Currency: matches[8],
		}
		posting.PriceType = matches[6]
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
		s += fmt.Sprintf(" %s %s", p.PriceType, p.Price)
	}
	s += "\n"

	return s
}
