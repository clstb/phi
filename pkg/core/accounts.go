package core

import (
	"context"

	"github.com/clstb/phi/pkg/core/db"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
)

func (s *Server) CreateAccount(
	ctx context.Context,
	req *pb.Account,
) (*pb.Account, error) {
	account := fin.AccountFromPB(req)

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
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	account = fin.NewAccount(accountDB)

	return account.PB(), nil
}

func (s *Server) GetAccounts(
	ctx context.Context,
	req *pb.AccountsQuery,
) (*pb.Accounts, error) {
	q := db.New(s.db)
	accountsDB, err := q.GetAccounts(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	accounts := fin.NewAccounts(accountsDB...)

	return accounts.PB(), nil
}
