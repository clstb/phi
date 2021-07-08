package fin

import (
	"fmt"
	"regexp"

	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

// Transactions is a slice of transaction
type Transactions []Transaction

// TransactionsFromPB creates a new transaction slice from it's protobuf representation.
func TransactionsFromPB(pb *pb.Transactions) (Transactions, error) {
	var transactions Transactions
	for _, v := range pb.Data {
		transaction, err := TransactionFromPB(v)
		if err != nil {
			return Transactions{}, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// TransactionsFromDB creates new transactions from their database representation.
func TransactionsFromDB(ts []db.Transaction) Transactions {
	var transactions Transactions
	for _, v := range ts {
		transactions = append(transactions, TransactionFromDB(v))
	}

	return transactions
}

// PB marshalls the transactions into their protobuf representation.
func (t Transactions) PB() *pb.Transactions {
	var data []*pb.Transaction
	for _, transaction := range t {
		data = append(data, transaction.PB())
	}

	return &pb.Transactions{
		Data: data,
	}
}

// Sum calculates the sum of all transactions grouped by accounts.
func (t Transactions) Sum() (map[string]Amounts, error) {
	sum := make(map[string]Amounts)
	for _, transaction := range t {
		weight, err := transaction.Weight()
		if err != nil {
			return nil, err
		}

		vFrom, ok := sum[transaction.From.String()]
		if !ok {
			vFrom = Amounts{weight.Neg()}
		} else {
			vFrom = append(vFrom, weight.Neg())
		}
		sum[transaction.From.String()] = vFrom

		vTo, ok := sum[transaction.To.String()]
		if !ok {
			vTo = Amounts{weight}
		} else {
			vTo = append(vTo, weight)
		}
		sum[transaction.To.String()] = vTo
	}

	var err error
	for k, v := range sum {
		sum[k], err = v.Sum()
		if err != nil {
			return nil, err
		}
	}

	return sum, nil
}

// Clear adds transactions to balance each income and expenses account to 0.
// The amounts are moved to Equity:Earnings:Current as well as Equity:Earnings:Previous depending on date.
func (t Transactions) Clear(accounts Accounts) (Transactions, error) {
	ec := accounts.ByName("Equity:Earnings:Current")
	if ec.Empty() {
		return Transactions{}, ErrNotFound{
			kind: "account",
			name: "Equity:Earnings:Current",
		}
	}

	re := regexp.MustCompile("^(Income|Expenses)")

	sum, err := t.Sum()
	if err != nil {
		return Transactions{}, err
	}

	accountsById := accounts.ById()
	var transactions Transactions
	for accountId, amounts := range sum {
		account, ok := accountsById[accountId]
		if !ok {
			return Transactions{}, fmt.Errorf("account not fount: %s", accountId)
		}
		if !re.MatchString(account.Name) {
			continue
		}

		for _, amount := range amounts {
			transaction := Transaction{}
			from, err := uuid.FromString(accountId)
			if err != nil {
				return Transactions{}, err
			}
			transaction.From = from
			transaction.To = ec.ID
			transaction.Units = amount
			transactions = append(transactions, transaction)
		}

	}

	return append(transactions, t...), nil
}

// ByDate groups transactions by their date.
func (t Transactions) ByDate() map[string]Transactions {
	byDate := map[string]Transactions{}
	for _, transaction := range t {
		date := transaction.Date.Format("2006-01-02")
		transactions, ok := byDate[date]
		if !ok {
			byDate[date] = Transactions{transaction}
		} else {
			byDate[date] = append(transactions, transaction)
		}
	}

	return byDate
}
