package cmd

import (
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/clstb/phi/pkg/ui"
	"github.com/urfave/cli/v2"
)

func Journal(ctx *cli.Context) error {
	core, err := Core(ctx)
	if err != nil {
		return err
	}

	from, to := ctx.String("from"), ctx.String("to")

	accountsPB, err := core.GetAccounts(
		ctx.Context,
		&pb.AccountsQuery{},
	)
	if err != nil {
		return err
	}

	accounts, err := fin.AccountsFromPB(accountsPB)
	if err != nil {
		return err
	}

	transactionPB, err := core.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			From: from,
			To:   to,
		},
	)
	if err != nil {
		return err
	}

	transactions, err := fin.TransactionsFromPB(transactionPB)
	if err != nil {
		return err
	}

	ui := ui.New(
		ctx.Context,
		transactions,
		accounts,
		core,
	)

	return ui.Run()
}
