package main

import (
	"filippo.io/age"
	"github.com/clstb/phi/go/pkg/config"
	"github.com/clstb/phi/go/pkg/crypto"
	"github.com/clstb/phi/go/pkg/parser"
	"github.com/clstb/phi/go/pkg/services/bookkeeper/pb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func Store(ctx *cli.Context) error {
	bookkeeperHost := ctx.String("bookkeeper-host")
	conn, err := grpc.Dial(
		bookkeeperHost,
		grpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	client := pb.NewBookkeeperClient(conn)

	ledger, err := parser.Load(ctx.String("ledger"))
	if err != nil {
		return err
	}

	configPath := ctx.String("config")
	config, err := config.Load(configPath)
	if err != nil {
		return err
	}

	identity, err := age.ParseX25519Identity(config.Identity)
	if err != nil {
		return err
	}

	dk, err := age.GenerateX25519Identity()
	if err != nil {
		return err
	}

	encryptedLedger, err := ledger.MarshalBytes(crypto.AgeEncoder(dk.Recipient()))
	if err != nil {
		return err
	}

	return nil

}
