package main

import (
	"github.com/clstb/phi/go/pkg/config"
	"github.com/clstb/phi/go/pkg/services/auth/pb"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func Login(ctx *cli.Context) error {
	authHost := ctx.String("auth-host")
	conn, err := grpc.Dial(
		authHost,
		grpc.WithInsecure(),
	)
	if err != nil {
		return err
	}

	client := pb.NewAuthClient(conn)

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

	jwt, err := client.Login(
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
