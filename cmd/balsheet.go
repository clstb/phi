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
	"google.golang.org/grpc"
)

func BalSheet(ctx *cli.Context) error {
	date := ctx.String("date")

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewCoreClient(conn)

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

	re := regexp.MustCompile("^Equity")
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
