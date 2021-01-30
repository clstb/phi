package fin

import (
	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

// Account is associated to multiple postings, an user and have a name.
type Account struct {
	db.Account
}

// AccountFromPB creates a new account from it's protobuf representation.
func AccountFromPB(a *pb.Account) (Account, error) {
	id, err := uuid.FromString(a.Id)
	if err != nil {
		return Account{}, err
	}

	account := Account{}
	account.ID = id
	account.Name = a.Name

	return account, nil
}

// AccountFromDB creates a new accoutn from it's database representation.
func AccountFromDB(a db.Account) Account {
	return Account{Account: a}
}

// PB marshalls the account into it's protobuf representation.
func (a Account) PB() *pb.Account {
	return &pb.Account{
		Id:   a.ID.String(),
		Name: a.Name,
	}
}

// Empty returns true if the account is initialized with default values.
func (a Account) Empty() bool {
	return a == Account{}
}
