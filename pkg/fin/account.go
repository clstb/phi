package fin

import (
	"github.com/clstb/phi/pkg/db"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

type Account struct {
	db.Account
}

func NewAccount(a db.Account) Account {
	return Account{a}
}

func AccountFromPB(a *pb.Account) Account {
	account := db.Account{
		ID:   uuid.FromStringOrNil(a.Id),
		Name: a.Name,
	}

	return NewAccount(account)
}

func (a Account) PB() *pb.Account {
	return &pb.Account{
		Id:   a.ID.String(),
		Name: a.Name,
	}
}
