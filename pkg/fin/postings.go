package fin

import (
	"fmt"

	"github.com/clstb/phi/pkg/core/db"
	"github.com/clstb/phi/pkg/pb"
)

type Postings []Posting

func NewPostings(postingsDB ...db.Posting) Postings {
	var postings Postings
	for _, p := range postingsDB {
		postings = append(postings, NewPosting(p))
	}

	return postings
}

func (p Postings) Sum() map[string]db.Amounts {
	sums := make(map[string]db.Amounts)
	for _, posting := range p {
		weight := posting.Weight()
		sum, ok := sums[posting.Account.String()]
		if !ok {
			sum = db.Amounts{weight}
		} else {
			sum = append(sum, weight)
		}
		sums[posting.Account.String()] = sum
	}

	return sums
}

func (p Postings) PB() []*pb.Posting {
	var postings []*pb.Posting
	for _, posting := range p {
		postings = append(postings, posting.PB())
	}

	return postings
}

func PostingsFromPB(pb *pb.Postings) (Postings, error) {
	var postings Postings
	for _, v := range pb.Data {
		posting, err := PostingFromPB(v)
		if err != nil {
			return Postings{}, fmt.Errorf("data: %w", err)
		}
		postings = append(postings, posting)
	}

	return postings, nil
}
