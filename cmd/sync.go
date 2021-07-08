package cmd

import (
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
)

func Sync(ctx *cli.Context) error {
	tinkGW, err := TinkGW(ctx)
	if err != nil {
		return err
	}

	_, err = tinkGW.Sync(
		ctx.Context,
		&pb.SyncReq{},
	)
	if err != nil {
		return err
	}

	return nil
}
