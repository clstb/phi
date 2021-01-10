package cmd

import (
	"fmt"

	gotree "github.com/DiSiqueira/GoTree"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func Journal(ctx *cli.Context) error {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewCoreClient(conn)

	transactionsRes, err := client.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			Fields: &pb.TransactionsQuery_Fields{
				Date:     true,
				Postings: true,
			},
			From: "-infinity",
			To:   "+infinity",
		},
	)

	transactions := make(map[string][]*pb.Transaction)

	for _, transaction := range transactionsRes.Data {
		transactions[transaction.Date] = append(
			transactions[transaction.Date],
			transaction,
		)
	}

	tree := gotree.New("Transactions")
	for k, v := range transactions {
		date := tree.Add(k)
		for _, transaction := range v {
			date.Add(transaction.Id)
		}
	}

	fmt.Println(tree.Print())

	return nil
}
