package core

import (
	"context"
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
	sub := uuid.FromStringOrNil(subStr)

	req.Id = uuid.Nil.String()
	account, err := fin.AccountFromPB(req)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := db.New(tx)

	accountDB, err := q.CreateAccount(ctx, account.Name)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	account = fin.AccountFromDB(accountDB)

	_, err = q.LinkAccount(ctx, db.LinkAccountParams{
		Account: account.ID,
		User:    sub,
	})
	if err != nil {
		tx.Rollback()
		return nil, err

	}

	if err := tx.Commit(); err != nil {
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
	sub := uuid.FromStringOrNil(subStr)

	q := db.New(s.db)
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
