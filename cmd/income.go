package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
	"google.golang.org/grpc"
)

func Income(ctx *cli.Context) error {
	from, to := ctx.String("from"), ctx.String("to")

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewCoreClient(conn)

	accounts, err := client.GetAccounts(
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

	postingsPB, err := client.GetPostings(
		ctx.Context,
		&pb.PostingsQuery{
			Fields: &pb.PostingsQuery_Fields{
				Account: true,
				Units:   true,
			},
			From:        from,
			To:          to,
			AccountName: "^(Income|Expenses)",
		},
	)
	if err != nil {
		return err
	}

	postings := fin.NewPostings()
	if err := postings.FromPB(postingsPB); err != nil {
		return err
	}

	sum := postings.Sum()
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

	var s string
	for _, amount := range sumByCurrency {
		s += "\t" + amount.String()
	}
	fmt.Fprintf(w, "\t\nNet Income:%s\n", s)

	return w.Flush()
}
