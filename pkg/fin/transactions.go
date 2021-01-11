package fin

import (
	"fmt"
	"regexp"

	"github.com/clstb/phi/pkg/pb"
)

type Transactions struct {
	Data []Transaction
	byId map[string]int32
}

func NewTransactions() *Transactions {
	return &Transactions{}
}

func (t Transactions) Sum() Sum {
	sum := make(Sum)
	for _, transaction := range t.Data {
		m := transaction.Sum()
		sum = sum.Add(m)
	}

	return sum
}

func (t Transactions) Clear(accounts Accounts) (Transactions, error) {
	sum := t.Sum()

	ec, ok := accounts.ByName("Equity:Earnings:Current")
	if !ok {
		return Transactions{}, fmt.Errorf("couldn't find account by name: Equity:Expenses:Current")
	}

	re := regexp.MustCompile("^(Income|Expenses)")
	var clears []Transaction
	for accountId, amounts := range sum {
		account, ok := accounts.ById(accountId)
		if !ok {
			return Transactions{}, fmt.Errorf("couldn't find account by id: %s", accountId)
		}
		if !re.MatchString(account.Name) {
			continue
		}
		for _, amount := range amounts {
			a := amount
			transaction := Transaction{
				Postings: Postings{
					Data: []Posting{
						{
							Account: ec.Id,
							Units:   a,
						},
						{
							Account: account.Id,
							Units:   a.Neg(),
						},
					},
				},
			}
			clears = append(clears, transaction)
		}
	}

	return Transactions{
		Data: append(t.Data, clears...),
	}, nil
}

func TransactionsFromPB(pb *pb.Transactions) (Transactions, error) {
	var data []Transaction
	byId := make(map[string]int32)
	var i int32
	for _, v := range pb.Data {
		transaction, err := TransactionFromPB(v)
		if err != nil {
			return Transactions{}, fmt.Errorf("data: %w", err)
		}
		data = append(data, transaction)
		byId[transaction.Id] = i
	}

	return Transactions{
		Data: data,
		byId: byId,
	}, nil
}

func (t Transactions) PB() (*pb.Transactions, error) {
	var data []*pb.Transaction
	byId := make(map[string]int32)
	var i int32
	for _, transaction := range t.Data {
		pb, err := transaction.PB()
		if err != nil {
			return nil, err
		}

		data = append(data, pb)
		byId[pb.Id] = i
	}

	return &pb.Transactions{
		Data: data,
		ById: byId,
	}, nil
}
