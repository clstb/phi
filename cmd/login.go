package cmd

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
)

func Login(ctx *cli.Context) error {
	auth, err := Auth(ctx)
	if err != nil {
		return err
	}

	name := ctx.String("name")
	password := ctx.String("password")
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

	fmt.Println(jwt)

	return nil
}
