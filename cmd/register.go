package cmd

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
)

func Register(ctx *cli.Context) error {
	core, err := Core(ctx)
	if err != nil {
		return err
	}

	auth, err := Auth(ctx)
	if err != nil {
		return err
	}

	name := ctx.String("name")
	password := ctx.String("password")
	user, err := auth.Register(
		ctx.Context,
		&pb.User{
			Name:     name,
			Password: password,
		},
	)
	if err != nil {
		return err
	}
	fmt.Println(user)

	accounts := []string{
		"Equity:OpeningBalances",
		"Equity:Earnings:Current",
		"Equity:Earnings:Previous",
	}

	fmt.Println("Welcome to phi!")
	for _, account := range accounts {
		fmt.Printf("Creating account %s\n", account)
		_, err := core.CreateAccount(
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
