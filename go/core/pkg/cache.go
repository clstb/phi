package pkg

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	"runtime/debug"
	"time"
)

// UserTokenCache sessionId -> sessionToken
var UserTokenCache = createCache()

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

func PutClientSessionTokenInCache(id string, token string) {
	err := UserTokenCache.Set(context.TODO(), id, token, &store.Options{Cost: 2})
	if err != nil {
		debug.PrintStack()
	}
}
