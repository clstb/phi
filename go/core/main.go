package main

import (
	"github.com/clstb/phi/go/core/pkg"
	"github.com/clstb/phi/go/core/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"log"
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
		log.Fatal(err)
	}
}

func Serve(ctx *cli.Context) error {

	authClient := auth.NewClient(ctx.String("oauthkeeper-uri"))
	server := pkg.NewServer(authClient)

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.POST("/api/login", server.DoLogin)
	router.POST("/api/register", server.DoRegister)

	router.POST("/api/link-tink", server.LinkBank)
	router.POST("/api/sync-ledger", server.SyncLedger)

	return router.Run("0.0.0.0:8081")
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
