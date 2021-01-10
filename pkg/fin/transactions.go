package fin

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
)

type Transactions struct {
	Data []*Transaction
	ById map[string]int32
}

func (t *Transactions) FromPB(pb *pb.Transactions) error {
	var data []*Transaction
	byId := make(map[string]int32)
	var i int32
	for _, v := range pb.Data {
		transaction := NewTransaction()
		if err := transaction.FromPB(v); err != nil {
			return fmt.Errorf("data: %w", err)
		}
		data = append(data, transaction)
		byId[transaction.Id] = i
	}

	t.Data = data
	t.ById = byId

	return nil
}

func (t *Transactions) PB() (*pb.Transactions, error) {
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
