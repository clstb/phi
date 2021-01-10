package fin

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
)

type Postings struct {
	Data []*Posting
	ById map[string]int32
}

func NewPostings() *Postings {
	return &Postings{}
}

func (p *Postings) Sum() Sum {
	sum := make(Sum)
	for _, posting := range p.Data {
		account, ok := sum[posting.Account]
		if !ok {
			sum[posting.Account] = make(map[string]*Amount)
			account = sum[posting.Account]
		}

		weight := posting.Weight()
		sum, ok := account[weight.Currency]
		if !ok {
			account[weight.Currency] = weight
		} else {
			account[weight.Currency] = weight.Add(sum)
		}
	}

	return sum
}

func (p *Postings) FromPB(pb *pb.Postings) error {
	var data []*Posting
	byId := make(map[string]int32)
	var i int32
	for _, v := range pb.Data {
		posting := NewPosting()
		if err := posting.FromPB(v); err != nil {
			return fmt.Errorf("data: %w", err)
		}
		data = append(data, posting)
		byId[posting.Id] = i
	}

	p.Data = data
	p.ById = byId

	return nil
}

func (p *Postings) PB() (*pb.Postings, error) {
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
