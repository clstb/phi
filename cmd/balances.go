package cmd

import (
	"os"
	"text/tabwriter"

	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
	"google.golang.org/grpc"
)

func Balances(ctx *cli.Context) error {
	date := ctx.String("date")

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
			To: date,
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
	tree.SetMetaValue("Balances")

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

	return w.Flush()
}
