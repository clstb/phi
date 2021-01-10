package cmd

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func CreateAccount(ctx *cli.Context) error {
	name := ctx.String("name")

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewCoreClient(conn)

	account, err := client.CreateAccount(
		ctx.Context,
		&pb.Account{
			Name: name,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(account)

	return nil
}
