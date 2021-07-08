package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
)

func Journal(ctx *cli.Context) error {
	core, err := Core(ctx)
	if err != nil {
		return err
	}

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
	accountsById := accounts.ById()

	from, to := ctx.String("from"), ctx.String("to")
	transactionsPB, err := core.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			From: from,
			To:   to,
		},
	)
	if err != nil {
		return err
	}

	transactions, err := fin.TransactionsFromPB(transactionsPB)
	if err != nil {
		return err
	}

	byDate := transactions.ByDate()
	var byDateKeys []string
	for k := range byDate {
		byDateKeys = append(byDateKeys, k)
	}
	sort.Slice(byDateKeys, func(i, j int) bool {
		return byDateKeys[i] < byDateKeys[j]
	})

	tree := treeprint.New()
	tree.SetMetaValue("Journal")
	tree.SetValue("")

	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)

	for _, date := range byDateKeys {
		transactions := byDate[date]
		tb := tree.AddMetaBranch(date, "\t\t\t\t\t\t")
		for _, transaction := range transactions {
			from, ok := accountsById[transaction.From.String()]
			if !ok {
				return fmt.Errorf("account not found: %s", transaction.From)
			}
			to, ok := accountsById[transaction.To.String()]
			if !ok {
				return fmt.Errorf("account not found: %s", transaction.To)
			}
			tb.AddMetaBranch(transaction.Entity, fmt.Sprintf("\t%s\t=>\t%s\t%s\t[%s]\t(%s)",
				from.Name,
				to.Name,
				transaction.Units.String(),
				transaction.Cost.String(),
				transaction.Price.String(),
			))
		}
	}

	_, err = w.Write(tree.Bytes())
	if err != nil {
		return err
	}

	return w.Flush()
}
