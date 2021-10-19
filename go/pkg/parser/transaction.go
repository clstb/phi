package parser

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var (
	TransactionRE = regexp.MustCompile(
		fmt.Sprintf(
			`^(%s) (txn|!|\*) "([^"]*)"(?:$| "([^"]*)")?$`,
			DateRE.String(),
		),
	)
)

type Transactions []Transaction

func (ts Transactions) FilterPayee(payee string) Transactions {
	var transactions Transactions
	for _, transaction := range ts {
		if transaction.Payee == payee {
			transactions = append(
				transactions,
				transaction,
			)
		}
	}

	return transactions
}

func (ts Transactions) ByTinkID() map[string]Transaction {
	m := make(map[string]Transaction)
	for _, t := range ts {
		for _, f := range t.Metadata {
			if f.Key == "tink_id" {
				m[f.Value] = t
			}
		}
	}

	return m
}

type Transaction struct {
	Index     int
	Date      string
	Type      string
	Payee     string
	Narration string
	Postings  []Posting
	Metadata  []MetadataField
}

func parseTransaction(
	s string,
	scanner *bufio.Scanner,
	ledger Ledger,
) Ledger {
	matches := TransactionRE.FindStringSubmatch(s)

	payee := ""
	if len(matches[4]) > 0 {
		payee = matches[3]
	}

	transaction := Transaction{
		Date:      matches[1],
		Type:      matches[2],
		Payee:     payee,
		Narration: matches[4],
	}

	for scanner.Scan() {
		s = scanner.Text()

		switch {
		case MetadataFieldRE.MatchString(s):
			transaction.Metadata = append(
				transaction.Metadata,
				parseMetadata(s),
			)
		case PostingRE.MatchString(s):
			posting := parsePosting(s)
			if posting.Units.IsZero() {
				for _, amount := range transaction.Weights() {
					if !amount.IsZero() {
						posting.Units = amount.Neg()
					}
				}
			}
			transaction.Postings = append(
				transaction.Postings,
				posting,
			)
		default:
			return append(
				append(ledger, transaction),
				parse(s, scanner, ledger)...,
			)
		}
	}

	return append(ledger, transaction)
}

func (t Transaction) String() string {
	s := fmt.Sprintf(
		"%s %s",
		t.Date,
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

func (t Transaction) Weights() map[string]Amount {
	m := map[string]Amount{}
	for _, posting := range t.Postings {
		w := posting.Weight()

		v, ok := m[w.Currency]
		if !ok {
			m[w.Currency] = w
		} else {
			m[w.Currency] = v.Add(w)
		}
	}

	return m
}

func (t Transaction) Balanced() bool {
	for _, amount := range t.Weights() {
		if amount.IsZero() {
			continue
		}

		return false
	}

	return true
}

func (t Transaction) Negative() bool {
	for _, amount := range t.Weights() {
		if amount.IsNegative() {
			return true
		}
	}

	return false
}

func (t Transaction) NormalizedPayee() []string {
	badRunes := map[rune]struct{}{
		'0': {},
		'1': {},
		'2': {},
		'3': {},
		'4': {},
		'5': {},
		'6': {},
		'7': {},
		'8': {},
		'9': {},
		'.': {},
		'*': {},
		'+': {},
		'-': {},
		'/': {},
		',': {},
	}

	normalized := []string{}

	s := strings.Split(t.Payee, " ")
	for _, v := range s {
		v = strings.ToLower(v)
		v = strings.TrimFunc(v, func(r rune) bool {
			if _, ok := badRunes[r]; ok {
				return true
			}
			return false
		})
		if v == "" {
			continue
		}

		normalized = append(normalized, v)
	}

	return normalized
}
