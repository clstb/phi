package internal

import (
	"github.com/clstb/phi/core/internal/auth"
	"github.com/clstb/phi/core/internal/handlers"
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

func NewServer(authClient *auth.AuthClient) *handlers.CoreServer {
	return &handlers.CoreServer{
		AuthClient:     authClient,
		Logger:         createLogger(),
		UserTokenCache: createCache(),
	}
}

func createLogger() *zap.SugaredLogger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	_logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	_sugar := _logger.Sugar()
	return _sugar
}

func createCache() *cache.Cache {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	ristrettoStore := store.NewRistretto(ristrettoCache, &store.Options{Expiration: 2 * time.Hour})
	cacheManager := cache.New(ristrettoStore)
	return cacheManager
}
