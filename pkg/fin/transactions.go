package fin

import (
	"fmt"
	"regexp"

	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

// Transactions is a slice of transaction
type Transactions []Transaction

// Sum calculates the sum of all transactions grouped by the account each posting belongs to.
func (t Transactions) Sum() (map[string]Amounts, error) {
	sums := make(map[string]Amounts)
	for _, transaction := range t {
		postingsSum, err := transaction.Postings.Sum()
		if err != nil {
			return nil, err
		}

		for accountId, amounts := range postingsSum {
			sum, ok := sums[accountId]
			if !ok {
				sum = amounts
			} else {
				sum = append(sum, amounts...)
			}
			sums[accountId] = sum
		}
	}
	for k, v := range sums {
		sum, err := v.Sum()
		if err != nil {
			return nil, err
		}
		sums[k] = sum
	}

	return sums, nil
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

	var transactions Transactions
	for accountId, amounts := range sum {
		// these accounts always exist so we don't check for empty
		account := accounts.ById(accountId)
		if !re.MatchString(account.Name) {
			continue
		}

		for _, amount := range amounts {
			transaction := Transaction{}

			posting := Posting{}
			posting.Account = ec.ID
			posting.Units = amount
			transaction.Postings = append(transaction.Postings, posting)

			posting = Posting{}
			posting.Account = uuid.FromStringOrNil(accountId) // TODO: this can break
			posting.Units = amount.Neg()
			transaction.Postings = append(transaction.Postings, posting)

			transactions = append(transactions, transaction)
		}

	}

	return append(transactions, t...), nil
}

// TransactionsFromPB creates a new transaction slice from it's protobuf representation.
func TransactionsFromPB(pb *pb.Transactions) (Transactions, error) {
	var transactions Transactions
	for _, v := range pb.Data {
		transaction, err := TransactionFromPB(v)
		if err != nil {
			return Transactions{}, fmt.Errorf("data: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
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
