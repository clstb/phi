package main

import (
	"fmt"
	"github.com/clstb/phi/core/internal/config"
	server2 "github.com/clstb/phi/core/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "TINK_GW_URI",
				EnvVars: []string{"TINK_GW_URI"},
				Value:   config.DefTinkGWAddr,
			},
			&cli.StringFlag{
				Name:    "LEDGER_URI",
				EnvVars: []string{"LEDGER_URI"},
				Value:   config.DefLedgerAddr,
			},
			&cli.StringFlag{
				Name:    "ORY_URI",
				EnvVars: []string{"ORY_URI"},
				Value:   config.DefOryUri,
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

	fmt.Println("------------------------")
	fmt.Printf("TINK_GW_URI => %s\n", ctx.String("TINK_GW_URI"))
	fmt.Printf("LEDGER_URI  => %s\n", ctx.String("LEDGER_URI"))
	fmt.Printf("ORY_URI     => %s\n", ctx.String("ORY_URI"))
	fmt.Println("------------------------")

	server := server2.NewServer(ctx.String("ORY_URI"), ctx.String("TINK_GW_URI"), ctx.String("LEDGER_URI"))
	router := gin.Default()

	router.Use(CORSMiddleware())
	router.POST("/api/login", server.DoLogin)
	router.POST("/api/register", server.DoRegister)

	router.POST("/api/auth/link", server.AuthCodeLink)
	router.POST("/api/auth/token", server.ExchangeToToken)
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
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
