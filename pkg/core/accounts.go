package core

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

func (s *Server) CreateAccount(
	ctx context.Context,
	req *pb.Account,
) (*pb.Account, error) {
	subStr, ok := ctx.Value("sub").(string)
	if !ok {
		return nil, fmt.Errorf("context: missing subject")
	}
	sub, err := uuid.FromString(subStr)
	if err != nil {
		return nil, err
	}

	req.Id = uuid.Nil.String()
	account, err := fin.AccountFromPB(req)
	if err != nil {
		return nil, err
	}

	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("context: missing transaction")
	}
	q := db.New(tx)

	accountDB, err := q.CreateAccount(ctx, account.Name)
	if err != nil {
		return nil, err
	}
	account = fin.AccountFromDB(accountDB)

	_, err = q.LinkAccount(ctx, db.LinkAccountParams{
		Account: account.ID,
		User:    sub,
	})
	if err != nil {
		return nil, err
	}

	return account.PB(), nil
}

func (s *Server) GetAccounts(
	ctx context.Context,
	req *pb.AccountsQuery,
) (*pb.Accounts, error) {
	subStr, ok := ctx.Value("sub").(string)
	if !ok {
		return nil, fmt.Errorf("context: missing subject")
	}
	sub, err := uuid.FromString(subStr)
	if err != nil {
		return nil, err
	}

	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("context: missing transaction")
	}
	q := db.New(tx)

	accountsDB, err := q.GetAccounts(ctx, db.GetAccountsParams{
		Name: req.Name,
		User: sub,
	})
	if err != nil {
		return nil, err
	}

	accounts := fin.AccountsFromDB(accountsDB...)

	return accounts.PB(), nil
}
