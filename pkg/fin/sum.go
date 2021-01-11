package fin

type SumCurrency map[string]Amount

func (s SumCurrency) Add(sum SumCurrency) SumCurrency {
	m := make(SumCurrency)

	for k, v := range s {
		m[k] = v
	}
	for k, v := range sum {
		amount, ok := m[k]
		if !ok {
			m[k] = v
		} else {
			m[k] = amount.Add(v)
		}
	}

	return m
}

func (s SumCurrency) Copy() SumCurrency {
	m := make(SumCurrency)

	for k, v := range s {
		m[k] = v
	}

	return m
}

type Sum map[string]SumCurrency

func (s Sum) Add(sum Sum) Sum {
	m := make(Sum)

	for k, v := range s {
		m[k] = v
	}
	for k, v := range sum {
		sc, ok := m[k]
		if !ok {
			m[k] = v.Copy()
		} else {
			m[k] = sc.Add(v)
		}
	}

	return m
}

func (s Sum) ByCurrency() SumCurrency {
	m := make(SumCurrency)

	for _, sc := range s {
		m = m.Add(sc)
	}

	return m
}
