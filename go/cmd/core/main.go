package main

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"os"
	"runtime/debug"
)

func main() {
	app := &cli.App{
		Action: Serve,
	}
	if err := app.Run(os.Args); err != nil {
		debug.PrintStack()
		sugar.Fatal(err)
	}
}

func Serve(ctx *cli.Context) error {
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.POST("/api/login", doLogin)
	router.POST("/api/register", doRegister)
	router.POST("/api/link-tink", getTinkLink)
	router.POST("/api/sync-ledger", SyncLedger)
	return router.Run("0.0.0.0:8099")
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
