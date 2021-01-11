package fin

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
)

type Postings struct {
	Data []Posting
	byId map[string]int32
}

func NewPostings() *Postings {
	return &Postings{}
}

func (p *Postings) Sum() Sum {
	sum := make(Sum)
	for _, posting := range p.Data {
		weight := posting.Weight()
		m := Sum{posting.Account: SumCurrency{weight.Currency: weight}}
		sum = sum.Add(m)
	}

	return sum
}

func PostingsFromPB(pb *pb.Postings) (Postings, error) {
	var data []Posting
	byId := make(map[string]int32)
	var i int32
	for _, v := range pb.Data {
		posting, err := PostingFromPB(v)
		if err != nil {
			return Postings{}, fmt.Errorf("data: %w", err)
		}
		data = append(data, posting)
		byId[posting.Id] = i
	}

	return Postings{
		Data: data,
		byId: byId,
	}, nil
}

func (p Postings) PB() (*pb.Postings, error) {
	var data []*pb.Posting
	byId := make(map[string]int32)
	var i int32
	for _, posting := range p.Data {
		pb, err := posting.PB()
		if err != nil {
			return nil, err
		}

		data = append(data, pb)
		byId[pb.Id] = i
	}

	return &pb.Postings{
		Data: data,
		ById: byId,
	}, nil
}
