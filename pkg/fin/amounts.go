package fin

type Amounts []Amount

func (a Amounts) Sum() (Amounts, error) {
	sums := make(map[string]Amount)
	for _, amount := range a {
		sum, ok := sums[amount.Currency]
		if !ok {
			sum = amount
		} else {
			v, err := sum.Add(amount)
			if err != nil {
				return Amounts{}, err
			}
			sum = v
		}
		sums[amount.Currency] = sum
	}

	var amounts Amounts
	for _, v := range sums {
		amounts = append(amounts, v)
	}

	return amounts, nil
}

func (a Amounts) ByCurrency(currency string) Amount {
	for _, amount := range a {
		if amount.Currency == currency {
			return amount
		}
	}

	return Amount{}
}

func (a Amounts) Currencies() []string {
	var currencies []string
	for _, amount := range a {
		currencies = append(currencies, amount.Currency)
	}

	return currencies
}
