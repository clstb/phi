package fin

import (
	"fmt"
	"regexp"

	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

type Transactions []Transaction

func (t Transactions) Sum() map[string]Amounts {
	sums := make(map[string]Amounts)
	for _, transaction := range t {
		for accountId, amounts := range transaction.Postings.Sum() {
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
		sums[k] = v.Sum()
	}

	return sums
}

func (t Transactions) Clear(accounts Accounts) (Transactions, error) {
	ec, ok := accounts.ByName("Equity:Earnings:Current")
	if !ok {
		return Transactions{}, fmt.Errorf("couldn't find account by name: Equity:Expenses:Current")
	}

	re := regexp.MustCompile("^(Income|Expenses)")

	var transactions Transactions
	for accountId, amounts := range t.Sum() {
		account, ok := accounts.ById(accountId)
		if !ok {
			return Transactions{}, fmt.Errorf("couldn't find account by id: %s", accountId)
		}
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

func (t Transactions) PB() *pb.Transactions {
	var data []*pb.Transaction
	for _, transaction := range t {
		data = append(data, transaction.PB())
	}

	return &pb.Transactions{
		Data: data,
	}
}
