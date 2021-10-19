package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/clstb/phi/go/pkg/crypto"
)

type Ledger []fmt.Stringer

func (l Ledger) Save(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	for _, v := range l {
		_, err := fmt.Fprint(f, v)
		if err != nil {
			return err
		}
	}

	return f.Close()
}

func parse(
	s string,
	scanner *bufio.Scanner,
	ledger Ledger,
) Ledger {
	switch {
	case TransactionRE.MatchString(s):
		return parseTransaction(s, scanner, ledger)
	case OpenRE.MatchString(s):
		return parseOpen(s, scanner, ledger)
	default:
		if scanner.Scan() {
			return append(
				append(ledger, stringer(s)),
				parse(scanner.Text(), scanner, ledger)...,
			)
		} else {
			return ledger
		}
	}
}

func Parse(r io.Reader) Ledger {
	scanner := bufio.NewScanner(r)
	scanner.Scan()

	ledger := parse(scanner.Text(), scanner, Ledger{})
	for i, v := range ledger {
		if t, ok := v.(Transaction); ok {
			t.Index = i
			ledger[i] = t
		}
	}

	return ledger
}

func Load(path string) (Ledger, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return Ledger{}, err
	}

	ledger := Parse(f)

	return ledger, f.Close()
}

func (l Ledger) Opens() Opens {
	var opens []Open
	for _, v := range l {
		if t, ok := v.(Open); ok {
			opens = append(opens, t)
		}
	}

	return opens
}

func (l Ledger) Transactions() Transactions {
	var transactions []Transaction
	for _, v := range l {
		if t, ok := v.(Transaction); ok {
			transactions = append(transactions, t)
		}
	}

	return transactions
}

func (l Ledger) MarshalBytes(transcoder crypto.Transcoder) ([]byte, error) {
	b := &bytes.Buffer{}
	for _, v := range l {
		_, err := fmt.Fprint(b, v)
		if err != nil {
			return nil, err
		}

	}

	return transcoder(b.Bytes())
}

func UnmarshalBytes(b []byte, transcoder crypto.Transcoder) (Ledger, error) {
	b, err := transcoder(b)
	if err != nil {
		return nil, err
	}

	return Parse(bytes.NewReader(b)), nil
}
