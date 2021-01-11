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

func BalSheet(ctx *cli.Context) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	date := ctx.String("date")

	accountsPB, err := client.GetAccounts(
		ctx.Context,
		&pb.AccountsQuery{
			Fields: &pb.AccountsQuery_Fields{
				Name: true,
			},
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
			From: "-infinity",
			To:   date,
		},
	)

	transactions, err := fin.TransactionsFromPB(transactionsPB)
	if err != nil {
		return err
	}

	cleared, err := transactions.Clear(accounts)
	if err != nil {
		return err
	}

	sum := cleared.Sum()
	sumByCurrency := sum.ByCurrency()

	tree := treeprint.New()
	tree.SetMetaValue("Balance Sheet")

	re := regexp.MustCompile("^(Equity|Liabilities|Assets)")
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	_, err = w.Write(renderTree(
		tree,
		accounts.FilterName(re),
		sum,
		sumByCurrency,
	))
	if err != nil {
		return err
	}

	re = regexp.MustCompile("^Equity")
	nw := make(fin.SumCurrency)
	for accountId, amounts := range sum {
		account, ok := accounts.ById(accountId)
		if !ok {
			continue
		}
		if re.MatchString(account.Name) {
			nw = nw.Add(amounts)
		}
	}

	var s string
	for _, amount := range nw {
		s += "\t" + amount.StringRaw()
	}
	fmt.Fprintf(w, "\t\nNet Worth:%s\n", s)

	return w.Flush()
}
