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
func AccountFromPB(a *pb.Account) Account {
	account := Account{}
	account.ID = uuid.FromStringOrNil(a.Id)
	account.Name = a.Name

	return account
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
