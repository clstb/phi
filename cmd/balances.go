package cmd

import (
	"os"
	"text/tabwriter"
	"time"

	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/golang/protobuf/ptypes"
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Balances(ctx *cli.Context) error {
	client, err := getClient(ctx)
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

	transactionsPB, err := client.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			Fields: &pb.TransactionFields{
				Date:     true,
				Postings: true,
			},
			From: &timestamppb.Timestamp{Seconds: 0, Nanos: 0},
			To:   dateProto,
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
	tree.SetMetaValue("Balances")

	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	_, err = w.Write(renderTree(
		tree,
		accounts,
		sum,
	))
	if err != nil {
		return err
	}

	return w.Flush()
}
