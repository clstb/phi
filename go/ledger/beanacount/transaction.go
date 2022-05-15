package beanacount

import (
	"fmt"
	"time"
)

func parseTransaction(s string) Transaction {
	matches := TransactionRE.FindStringSubmatch(s)

	// ignore error as we verify through regexp
	date, _ := time.Parse("2006-01-02", matches[1])

	return Transaction{
		Date:      date,
		Type:      matches[2],
		Payee:     matches[3],
		Narration: matches[4],
	}
}

func (t Transaction) String() string {
	s := fmt.Sprintf(
		"%s %s",
		t.Date.Format("2006-01-02"),
		t.Type,
	)
	if t.Payee != "" {
		s = fmt.Sprintf(`%s "%s"`, s, t.Payee)
	}
	s = fmt.Sprintf(`%s "%s"`, s, t.Narration) + "\n"

	for _, md := range t.Metadata {
		s += md.String()
	}

	for _, posting := range t.Postings {
		s += posting.String()
	}

	return s
}

type Transactions []Transaction

func (t Transactions) ByTinkId() map[string]Transaction {
	m := map[string]Transaction{}
	for _, transaction := range t {
		for _, md := range transaction.Metadata {
			if md.Key == "tink_id" {
				m[md.Value[1:len(md.Value)-1]] = transaction
			}
		}
	}
	return m
}

func (t Transactions) ByPayee() map[string]Transactions {
	m := map[string]Transactions{}
	for _, transaction := range t {
		m[transaction.Payee] = append(m[transaction.Payee], transaction)
	}
	return m
}

func (t Transactions) ByMonth() map[string]Transactions {
	m := map[string]Transactions{}
	for _, transaction := range t {
		s := transaction.Date.Format("2006-01")
		m[s] = append(m[s], transaction)
	}
	return m
}

func (t Transactions) AccountsByPayee() map[string][]string {
	m := map[string][]string{}
	for _, transaction := range t {
		var accounts []string
		for _, posting := range transaction.Postings {
			accounts = append(accounts, posting.Account)
		}
		m[transaction.Payee] = append(m[transaction.Payee], accounts...)
	}

	for k, v := range m {
		keys := map[string]bool{}
		var unique []string
		for _, s := range v {
			if _, ok := keys[s]; !ok {
				keys[s] = true
				unique = append(unique, s)
			}
		}
		m[k] = unique
	}

	return m
}
