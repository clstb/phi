package main

import (
	"github.com/clstb/phi/go/pkg/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"os"

	"github.com/clstb/phi/go/internal/tinkgw/server"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "tink-client-id",
				EnvVars:  []string{"TINK_CLIENT_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "tink-client-secret",
				EnvVars:  []string{"TINK_CLIENT_SECRET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "ory-token",
				EnvVars:  []string{"ORY_TOKEN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "jwks-url",
				EnvVars:  []string{"JWKS_URL"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "callback-url",
				EnvVars:  []string{"CALLBACK_URL"},
				Required: true,
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 8080,
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	server, err := server.NewServer(
		ctx.String("tink-client-id"),
		ctx.String("tink-client-secret"),
		ctx.String("ory-token"),
		ctx.String("callback-url"),
	)
	if err != nil {
		return err
	}

	router := gin.Default()
	router.Use(
		middleware.Auth(server.Logger, ctx.String("jwks-url")),
	)

	router.GET("/api/link", server.Link)
	router.POST("/api/token", server.Token)
	router.POST("/api/tink-user", func(context *gin.Context) {
		server.RegisterTinkUser(server.OryToken, context)
	})
	return router.Run("0.0.0.0:8080")
}
