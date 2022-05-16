package beanacount

import (
	"fmt"
)

func (p Posting) String() string {
	if p.Units.IsZero() {
		return fmt.Sprintf("  %s\n", p.Account.String())
	}

	s := fmt.Sprintf(
		"  %s %s",
		p.Account.String(),
		p.Units,
	)
	if !p.Cost.IsZero() {
		s += fmt.Sprintf(" {%s}", p.Cost.String())
	}
	if !p.Price.IsZero() {
		s += fmt.Sprintf(" %s %s", p.PriceType, p.Price.String())
	}
	s += "\n"

	return s
}
