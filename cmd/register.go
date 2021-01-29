package cmd

import (
	"fmt"

	"github.com/clstb/phi/pkg/config"
	"github.com/clstb/phi/pkg/pb"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func Register(ctx *cli.Context) error {
	p := promptui.Prompt{
		Label: "Username",
	}

	name, err := p.Run()
	if err != nil {
		return err
	}

	p = promptui.Prompt{
		Label: "Password",
		Mask:  '*',
	}
	password, err := p.Run()
	if err != nil {
		return err
	}

	p = promptui.Prompt{
		Label: "Retype Password",
		Mask:  '*',
	}
	retyped, err := p.Run()
	if err != nil {
		return err
	}

	if password != retyped {
		return fmt.Errorf("passwords don't match")
	}

	auth, err := Auth(ctx)
	if err != nil {
		return err
	}

	jwt, err := auth.Register(
		ctx.Context,
		&pb.User{
			Name:     name,
			Password: password,
		},
	)
	if err != nil {
		return err
	}

	configPath := ctx.String("config")
	config, err := config.Load(configPath)
	if err != nil {
		return err
	}
	config.AccessToken = jwt.AccessToken

	if err := config.Save(configPath); err != nil {
		return err
	}

	core, err := Core(ctx)
	if err != nil {
		return err
	}

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
