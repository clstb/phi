package fin_test

import (
	"testing"

	"github.com/clstb/phi/pkg/fin"
	"github.com/matryer/is"
)

func TestAmountsSum(t *testing.T) {
	is := is.New(t)

	type test struct {
		do    func() (fin.Amounts, error)
		check func(fin.Amounts, error)
	}
	var tests []test
	add := func(t test) {
		tests = append(tests, t)
	}
	add(test{
		do: func() (fin.Amounts, error) {
			return fin.Amounts{
				fin.NewAmount(1, 1, "EUR"),
				fin.NewAmount(1, 1, "EUR"),
				fin.NewAmount(1, 1, "EUR"),
				fin.NewAmount(1, 1, "EUR"),
				fin.NewAmount(1, 1, "EUR"),
				fin.NewAmount(3, 1, "USD"),
				fin.NewAmount(3, 1, "USD"),
				fin.NewAmount(3, 1, "USD"),
				fin.NewAmount(3, 1, "USD"),
				fin.NewAmount(3, 1, "USD"),
				fin.NewAmount(1, 0, "CAD"),
				fin.NewAmount(1, 0, "CAD"),
				fin.NewAmount(1, 0, "CAD"),
				fin.NewAmount(1, 0, "CAD"),
				fin.NewAmount(1, 0, "CAD"),
			}.Sum()
		},
		check: func(a fin.Amounts, e error) {
			is.NoErr(e)
			is.Equal(a.ByCurrency("EUR"), fin.NewAmount(5, 1, "EUR"))
			is.Equal(a.ByCurrency("USD"), fin.NewAmount(15, 1, "USD"))
			is.Equal(a.ByCurrency("CAD"), fin.NewAmount(5, 0, "CAD"))
		},
	})

	for _, t := range tests {
		t.check(t.do())
	}
}
