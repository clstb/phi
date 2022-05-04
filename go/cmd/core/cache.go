package main

import (
	"context"
	"github.com/clstb/phi/go/pkg/client"
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	"time"
)

// UserClientCache sessionId -> oryClient
var UserClientCache = createCache()

// LoggedInUsersCache username -> sessionId
var LoggedInUsersCache = createCache()

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

func putSessionIdInCache(name string, id string) {
	err := LoggedInUsersCache.Set(context.TODO(), name, id, &store.Options{Cost: 2})
	if err != nil {
		sugar.Error(err)
	}
}

func putClientInCache(id string, c *client.Client) {
	err := UserClientCache.Set(context.TODO(), id, c, &store.Options{Cost: 2})
	if err != nil {
		sugar.Error(err)
	}
}
