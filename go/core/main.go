package main

import (
	"github.com/clstb/phi/go/tinkgw/pkg"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"os"
	"runtime/debug"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "oauthkeeper-uri",
				EnvVars:  []string{"OAUTHKEEPER_URI"},
				Required: true,
			},
		},
		Action: Serve,
	}

	if err := app.Run(os.Args); err != nil {
		debug.PrintStack()
		pkg.Sugar.Fatal(err)
	}
}

func Serve(ctx *cli.Context) error {
	/*
		authClient := client.NewClient(OauthKeeperUri)

		router := gin.Default()
		router.Use(CORSMiddleware())
		router.POST("/api/login", func(context *gin.Context) {
			handlers.DoLogin(context, authClient)
		})
		router.POST("/api/register", func(context *gin.Context) {
			handlers.DoRegister(context, authClient)
		})
		router.POST("/api/link-tink", func(context *gin.Context) {
			handlers.GetTinkLink(context, authClient)
		})
		router.POST("/api/sync-ledger", func(context *gin.Context) {
			handlers.SyncLedger(context, authClient)
		})
		return router.Run("0.0.0.0:" + core.ServerPort)
	*/
	return nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
