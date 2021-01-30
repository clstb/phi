package fin

import (
	"regexp"

	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
)

// Accounts is a slice of account
type Accounts []Account

// AccountsFromPB creates a new account slice from it's protobuf representation.
func AccountsFromPB(a *pb.Accounts) (Accounts, error) {
	var accounts Accounts
	for _, v := range a.Data {
		account, err := AccountFromPB(v)
		if err != nil {
			return Accounts{}, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// AccountsFromDB create a new account slice from it's database representation.
func AccountsFromDB(a ...db.Account) Accounts {
	var accounts Accounts
	for _, account := range a {
		accounts = append(accounts, AccountFromDB(account))
	}

	return accounts
}

// PB marshalls the account into it's protobuf representation.
func (a Accounts) PB() *pb.Accounts {
	var data []*pb.Account
	for _, account := range a {
		data = append(data, account.PB())
	}

	return &pb.Accounts{
		Data: data,
	}
}

// ById returns the first account matching given id.
// An empty account is returned when no account is found.
func (a Accounts) ById(id string) Account {
	for _, account := range a {
		if account.ID.String() == id {
			return account
		}
	}
	return Account{}
}

// ByName returns the first account matching given name.
// An empty account is returned when no account is found.
func (a Accounts) ByName(name string) Account {
	for _, account := range a {
		if account.Name == name {
			return account
		}
	}
	return Account{}
}

// Names returns a string slice of account names.
func (a Accounts) Names() []string {
	var names []string
	for _, account := range a {
		names = append(names, account.Name)
	}

	return names
}

// FilterName returns a string slice of account names filtered by provided regexp.
// It checks whether or not the regexp matches for each account name.
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
