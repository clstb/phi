package main

import (
	"github.com/clstb/phi/go/core/pkg/client"
	"github.com/clstb/phi/go/core/pkg/handlers"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"runtime/debug"
	"time"
)

var _, sugar = func() (*zap.Logger, *zap.SugaredLogger) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	_logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	_sugar := _logger.Sugar()
	return _logger, _sugar
}()

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
		sugar.Fatal(err)
	}
}

func Serve(ctx *cli.Context) error {

	authClient := client.NewClient(ctx.String("oauthkeeper-uri"))
	server := handlers.CoreServer{AuthClient: authClient, Logger: sugar}

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.POST("/api/login", server.DoLogin)
	router.POST("/api/register", server.DoRegister)

	/*
		router.POST("/api/link-tink", func(context *gin.Context) {
			handlers.GetTinkLink(context, authClient)
		})
		router.POST("/api/sync-ledger", func(context *gin.Context) {
			handlers.SyncLedger(context, authClient)
		})

	*/
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
