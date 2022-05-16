package main

import (
	"github.com/clstb/phi/core/internal/auth"
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
		Action: Serve,
	}

	if err := app.Run(os.Args); err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}
}

func Serve(ctx *cli.Context) error {

	authClient := auth.NewClient()
	server := server2.NewServer(authClient)

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
