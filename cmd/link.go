package cmd

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
)

func Link(ctx *cli.Context) error {
	tinkGW, err := TinkGW(ctx)
	if err != nil {
		return err
	}

	link, err := tinkGW.Link(
		ctx.Context,
		&pb.LinkReq{
			Market: "DE",
			Locale: "de_DE",
		},
	)
	if err != nil {
		return err
	}
	fmt.Println(link.TinkLink)

	return nil
}
