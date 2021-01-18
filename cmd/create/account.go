package create

import (
	"fmt"

	"github.com/clstb/phi/cmd"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
)

func Account(ctx *cli.Context) error {
	core, err := cmd.Core(ctx)
	if err != nil {
		return err
	}

	name := ctx.String("name")

	account, err := core.CreateAccount(
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
