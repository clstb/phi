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

	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)

	for _, date := range byDateKeys {
		transactions := byDate[date]
		tb := tree.AddMetaBranch(date, "\t")
		for _, transaction := range transactions {
			pb := tb.AddMetaBranch(transaction.Entity, "\t")
			for _, posting := range transaction.Postings {
				// these accounts always exist so we don't check for empty
				account := accounts.ById(posting.Account.String())
				pb.AddMetaNode(account.Name, fmt.Sprintf(
					"\t%s\t%s\t%s",
					posting.Units.Color(false),
					posting.Cost.Color(false),
					posting.Price.Color(false),
				))
			}
		}
	}

	_, err = w.Write(tree.Bytes())
	if err != nil {
		return err
	}

	return w.Flush()
}
