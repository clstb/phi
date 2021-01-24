package fin_test

import (
	"testing"

	"github.com/clstb/phi/pkg/fin"
	"github.com/matryer/is"
)

func TestAmounts(t *testing.T) {
	is := is.New(t)

	type test struct {
		do    func() fin.Amounts
		check func(fin.Amounts)
	}
	var tests []test
	add := func(t test) {
		tests = append(tests, t)
	}
	fromString := func(s string) fin.Amount {
		amount, err := fin.AmountFromString(s)
		is.NoErr(err)
		return amount
	}
	add(test{
		do: func() fin.Amounts {
			return fin.Amounts{
				fromString("0.1 EUR"),
				fromString("0.1 EUR"),
				fromString("0.1 EUR"),
				fromString("0.1 EUR"),
				fromString("0.1 EUR"),
				fromString("0.3 USD"),
				fromString("0.3 USD"),
				fromString("0.3 USD"),
				fromString("0.3 USD"),
				fromString("0.3 USD"),
				fromString("1 CAD"),
				fromString("1 CAD"),
				fromString("1 CAD"),
				fromString("1 CAD"),
				fromString("1 CAD"),
			}.Sum()
		},
		check: func(a fin.Amounts) {
			is.Equal(a.ByCurrency("EUR"), fromString("0.5 EUR"))
			is.Equal(a.ByCurrency("USD"), fromString("1.5 USD"))
			is.Equal(a.ByCurrency("CAD"), fromString("5 CAD"))
			is.Equal(a.Currencies(), []string{"EUR", "USD", "CAD"})
		},
	})

	for _, t := range tests {
		t.check(t.do())
	}
}
