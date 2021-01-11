package fin

import "github.com/clstb/phi/pkg/pb"

type Account struct {
	Id   string
	Name string
}

func AccountFromPB(pb *pb.Account) Account {
	return Account{
		Id:   pb.Id,
		Name: pb.Name,
	}

}

func (a Account) PB() *pb.Account {
	return &pb.Account{
		Id:   a.Id,
		Name: a.Id,
	}
}
