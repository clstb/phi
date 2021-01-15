package cmd

import (
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"
	"time"

	"github.com/clstb/phi/pkg/db"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
)

func BalSheet(ctx *cli.Context) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	accountsPB, err := client.GetAccounts(
		ctx.Context,
		&pb.AccountsQuery{
			Fields: &pb.AccountFields{
				Name: true,
			},
		},
	)
	if err != nil {
		return err
	}
	accounts, err := fin.AccountsFromPB(accountsPB)
	if err != nil {
		return err
	}

	date, err := time.Parse("2006-01-02", ctx.String("date"))
	if err != nil {
		return err
	}
	dateProto, err := ptypes.TimestampProto(date)
	if err != nil {
		return err
	}

	transactionsPB, err := client.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			Fields: &pb.TransactionFields{
				Date:     true,
				Postings: true,
			},
			From: &timestamp.Timestamp{Seconds: 0, Nanos: 0},
			To:   dateProto,
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

	tree := treeprint.New()
	tree.SetMetaValue("Balance Sheet")

	re := regexp.MustCompile("^(Equity|Liabilities|Assets)")
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	_, err = w.Write(renderTree(
		tree,
		accounts.FilterName(re),
		sum,
	))
	if err != nil {
		return err
	}

	re = regexp.MustCompile("^Equity")
	var amounts db.Amounts
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

	nw := amounts.Sum()
	var s string
	for _, amount := range nw {
		s += "\t" + amount.StringRaw()
	}
	fmt.Fprintf(w, "\t\nNet Worth:%s\n", s)

	return w.Flush()
}
