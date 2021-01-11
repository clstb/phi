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
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	from, to := ctx.String("from"), ctx.String("to")

	accountsPB, err := client.GetAccounts(
		ctx.Context,
		&pb.AccountsQuery{
			Fields: &pb.AccountsQuery_Fields{
				Name: true,
			},
			Name: "^(Income|Expenses)",
		},
	)
	if err != nil {
		return err
	}
	accounts := fin.AccountsFromPB(accountsPB)

	transactionsPB, err := client.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			Fields: &pb.TransactionsQuery_Fields{
				Date:     true,
				Postings: true,
			},
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

	sum := transactions.Sum()

	sumByCurrency := sum.ByCurrency()

	tree := treeprint.New()
	tree.SetMetaValue("Income Statement")

	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	_, err = w.Write(renderTree(
		tree,
		accounts,
		sum,
		sumByCurrency,
	))
	if err != nil {
		return err
	}

	re := regexp.MustCompile("^(Income|Expenses)")
	ni := make(fin.SumCurrency)
	for accountId, amounts := range sum {
		account, ok := accounts.ById(accountId)
		if !ok {
			continue
		}
		if re.MatchString(account.Name) {
			ni = ni.Add(amounts)
		}
	}

	var s string
	for _, amount := range ni {
		s += "\t" + amount.StringRaw()
	}
	fmt.Fprintf(w, "\t\nNet Income:%s\n", s)

	return w.Flush()
}
