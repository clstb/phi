package cmd

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func Register(ctx *cli.Context) error {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewCoreClient(conn)

	accounts := []string{
		"Equity:OpeningBalances",
		"Equity:Earnings:Current",
		"Equity:Earnings:Previous",
	}

	fmt.Println("Welcome to phi!")
	for _, account := range accounts {
		fmt.Printf("Creating account %s\n", account)
		_, err := client.CreateAccount(
			ctx.Context,
			&pb.Account{
				Name: account,
			},
		)
		if err != nil {
			return err
		}
	}
	fmt.Println("Success!")

	return nil
}
