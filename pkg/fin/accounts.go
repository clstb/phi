package fin

import (
	"regexp"

	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
)

type Accounts []Account

func NewAccounts(accountsDB ...db.Account) Accounts {
	var accounts Accounts
	for _, a := range accountsDB {
		accounts = append(accounts, NewAccount(a))
	}

	return accounts
}

func AccountsFromPB(pb *pb.Accounts) (Accounts, error) {
	var accounts Accounts
	for _, v := range pb.Data {
		accounts = append(accounts, AccountFromPB(v))
	}

	return accounts, nil
}

func (a Accounts) PB() *pb.Accounts {
	var data []*pb.Account
	for _, account := range a {
		data = append(data, account.PB())
	}

	return &pb.Accounts{
		Data: data,
	}
}

func (a Accounts) ById(id string) (Account, bool) {
	for _, account := range a {
		if account.ID.String() == id {
			return account, true
		}
	}
	return Account{}, false
}

func (a Accounts) ByName(name string) (Account, bool) {
	for _, account := range a {
		if account.Name == name {
			return account, true
		}
	}
	return Account{}, false
}

func (a Accounts) Names() []string {
	var names []string
	for _, account := range a {
		names = append(names, account.Name)
	}

	return names
}

func (a Accounts) FilterName(re *regexp.Regexp) Accounts {
	var accounts Accounts

	for _, account := range a {
		if !re.MatchString(account.Name) {
			continue
		}
		accounts = append(accounts, account)
	}

	return accounts
}
