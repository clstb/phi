package cmd

import (
	"github.com/clstb/phi/pkg/config"
	"github.com/clstb/phi/pkg/pb"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func Login(ctx *cli.Context) error {
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

	auth, err := Auth(ctx)
	if err != nil {
		return err
	}

	jwt, err := auth.Login(
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

	return nil
}
