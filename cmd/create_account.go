package cmd

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
)

func CreateAccount(ctx *cli.Context) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	name := ctx.String("name")

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
