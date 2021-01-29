package cmd

import (
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"

	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
)

func Income(ctx *cli.Context) error {
	core, err := Core(ctx)
	if err != nil {
		return err
	}

	accountsPB, err := core.GetAccounts(
		ctx.Context,
		&pb.AccountsQuery{
			Name: "^(Income|Expenses)",
		},
	)
	if err != nil {
		return err
	}
	accounts, err := fin.AccountsFromPB(accountsPB)
	if err != nil {
		return err
	}

	from, to := ctx.String("from"), ctx.String("to")
	transactionsPB, err := core.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			From:        from,
			To:          to,
			AccountName: "^(Income|Expenses)",
		},
	)
	if err != nil {
		return err
	}

	transactions, err := fin.TransactionsFromPB(transactionsPB)
	if err != nil {
		return err
	}

	sum, err := transactions.Sum()
	if err != nil {
		return err
	}

	amounts := fin.Amounts{}
	for _, v := range sum {
		amounts = append(amounts, v...)
	}
	amountsSum, err := amounts.Sum()
	if err != nil {
		return err
	}
	currencies := amountsSum.Currencies()

	tree := treeprint.New()
	tree.SetMetaValue("Income Statement")

	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	_, err = w.Write(renderTree(
		tree,
		accounts,
		currencies,
		sum,
	))
	if err != nil {
		return err
	}

	re := regexp.MustCompile("^(Income|Expenses)")
	amounts = fin.Amounts{}
	for accountId, v := range sum {
		account, ok := accounts.ById(accountId)
		if !ok {
			continue
		}
		if !re.MatchString(account.Name) {
			continue
		}
		amounts = append(amounts, v...)
	}

	ni, err := amounts.Sum()
	if err != nil {
		return err
	}

	var s string
	for _, amount := range ni {
		s += "\t" + amount.ColorRaw(true)
	}
	fmt.Fprintf(w, "\t\nNet Income:%s\n", s)

	return w.Flush()
}
