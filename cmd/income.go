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
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
)

func Income(ctx *cli.Context) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	accountsPB, err := client.GetAccounts(
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

	from, err := time.Parse("2006-01-02", ctx.String("from"))
	if err != nil {
		return err
	}
	fromProto, err := ptypes.TimestampProto(from)
	if err != nil {
		return err
	}
	to, err := time.Parse("2006-01-02", ctx.String("to"))
	if err != nil {
		return err
	}
	toProto, err := ptypes.TimestampProto(to)
	if err != nil {
		return err
	}

	transactionsPB, err := client.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			From:        fromProto,
			To:          toProto,
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

	tree := treeprint.New()
	tree.SetMetaValue("Income Statement")

	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	_, err = w.Write(renderTree(
		tree,
		accounts,
		sum,
	))
	if err != nil {
		return err
	}

	re := regexp.MustCompile("^(Income|Expenses)")
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

	ni := amounts.Sum()
	var s string
	for _, amount := range ni {
		s += "\t" + amount.StringRaw()
	}
	fmt.Fprintf(w, "\t\nNet Income:%s\n", s)

	return w.Flush()
}
